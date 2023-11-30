package service

import (
	"axon/internal/integrations"
	"axon/internal/integrations/dto"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/repository"
	"axon/utils"
	zkErrorsAxon "axon/utils/zkerrors"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/model"
	"github.com/zerok-ai/zk-utils-go/common"
	zkUtils "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/ds"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var LogTag = "zk_prometheus_service"

type PrometheusService interface {
	GetPodsInfoService(podInfoReq request.PromRequestMeta) (promResponse.PodsInfoResponse, *zkerrors.ZkError)
	GetContainerInfoService(podInfoReq request.PromRequestMeta) (promResponse.ContainerInfoResponse, *zkerrors.ZkError)
	GetContainerMetricService(podInfoReq request.PromRequestMeta) (promResponse.ContainerMetricsResponse, *zkerrors.ZkError)
	GetGenericQueryService(genericQueryReq request.GenericPromRequest) (promResponse.GenericQueryResponse, *zkerrors.ZkError)
	TestIntegrationConnection(integrationId string) (promResponse.TestConnectionResponse, *zkerrors.ZkError)
	TestUnsavedIntegrationConnection(url string, username, password *string) (promResponse.TestConnectionResponse, *zkerrors.ZkError)
	GetMetricAttributes(integrationId string, metricName string, st string, et string) (promResponse.MetricAttributesListResponse, *zkerrors.ZkError)
	MetricsList(integrationId string) (promResponse.IntegrationMetricsListResponse, *zkerrors.ZkError)
	AlertsList(integrationId string, name string) (promResponse.IntegrationAlertsListResponse, *zkerrors.ZkError)
	GetAlertsTimeSeriesTrigger(integrationId string, step string, time string, endTime string) (promResponse.AlertTimeSeriesResponse, *zkerrors.ZkError)
	GetMetricServerRepo() repository.PromQLRepo
	PrometheusAlertWebhook(string, promResponse.AlertWebhookResponse)
}

func NewPrometheusService(integrationsManager integrations.IntegrationsManager) PrometheusService {
	promIntegrations := make(map[string]repository.PromQLRepo)
	return prometheusService{integrationsManager: integrationsManager, promIntegrations: promIntegrations}
}

type prometheusService struct {
	integrationsManager integrations.IntegrationsManager
	metricServerRepo    repository.PromQLRepo
	metricServerId      string
	promIntegrations    map[string]repository.PromQLRepo
}

func (s prometheusService) GetMetricServerRepo() repository.PromQLRepo {
	promIntegrations := s.integrationsManager.GetIntegrationsByType(dto.PrometheusIntegrationType)
	for _, promIntegration := range promIntegrations {
		if promIntegration.Disabled == false && promIntegration.Deleted == false && promIntegration.MetricServer {
			if s.metricServerId == promIntegration.Id {
				return s.metricServerRepo
			}
			// Create new client for metric server
			promClient, err := api.NewClient(api.Config{
				Address: promIntegration.Url,
			})
			if err != nil {
				zkLogger.Error(LogTag, "Error creating Prometheus client: %v\n", err)
				continue
			}
			s.metricServerRepo = repository.NewPromQLRepo(promClient)
			s.metricServerId = promIntegration.Id
			return s.metricServerRepo
		}
	}
	zkLogger.Error(LogTag, "No metric server found")
	return nil
}

func (s prometheusService) GetPromIntegrationById(integrationId string) repository.PromQLRepo {
	integrationItem := s.integrationsManager.GetIntegrationById(integrationId)
	if integrationItem == nil || integrationItem.Type != dto.PrometheusIntegrationType || integrationItem.Disabled == true || integrationItem.Deleted == true {
		zkLogger.Error(LogTag, "Missing integration id: ", integrationId)
		return nil
	}
	if s.promIntegrations[integrationId] != nil {
		return s.promIntegrations[integrationId]
	}
	// Create integration client
	promClient, err := api.NewClient(api.Config{
		Address: integrationItem.Url,
	})
	if err != nil {
		zkLogger.Error(LogTag, "Error creating Prometheus client: %v\n", err)
		return nil
	}
	return repository.NewPromQLRepo(promClient)
}

func (s prometheusService) GetGenericQueryService(genericQueryReq request.GenericPromRequest) (promResponse.GenericQueryResponse, *zkerrors.ZkError) {
	var response promResponse.GenericQueryResponse
	targetPromIntegration := s.GetPromIntegrationById(genericQueryReq.PromIntegrationId)
	if targetPromIntegration == nil {
		respErr := utils.BuildZkError(LogTag, "No prometheus integration found for id: ", genericQueryReq.PromIntegrationId)
		return response, respErr
	}
	queryResult, resultType, err := targetPromIntegration.GenericQuery(genericQueryReq)
	if err != nil {
		respErr := utils.BuildZkError(LogTag, "Failed to query prometheus, Error: ", err.Error())
		return response, respErr
	}
	response.Result = queryResult
	response.Type = resultType
	return response, nil
}

func (s prometheusService) GetPodsInfoService(podInfoReq request.PromRequestMeta) (promResponse.PodsInfoResponse, *zkerrors.ZkError) {
	var response promResponse.PodsInfoResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		zkLogger.Error(LogTag, "No metric server found")
		respErr := utils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podsInfo, err := metricServerRepo.PodsInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podInfo: ", err.Error())
		respErr := utils.BuildZkError(LogTag, "Error while collecting podInfo")
		return response, respErr
	}

	podsInfoItems := extractMetricAttributes(podsInfo)
	response.PodsInfo = podsInfoItems
	return response, nil
}

func (s prometheusService) GetContainerInfoService(podInfoReq request.PromRequestMeta) (promResponse.ContainerInfoResponse, *zkerrors.ZkError) {
	var response promResponse.ContainerInfoResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		respErr := utils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podContainerInfo, err := metricServerRepo.PodContainerInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podContainerInfo: ", err.Error())
		respErr := utils.BuildZkError(LogTag, "Error while collecting podContainerInfo")
		return response, respErr
	}
	podContainerInfoItems := extractMetricAttributes(podContainerInfo)
	response.ContainerInfo = podContainerInfoItems
	return response, nil
}

func (s prometheusService) GetContainerMetricService(podInfoReq request.PromRequestMeta) (promResponse.ContainerMetricsResponse, *zkerrors.ZkError) {
	var response promResponse.ContainerMetricsResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		respErr := utils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}

	cpuUsageData, err := metricServerRepo.GetPodCPUUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		respErr := utils.BuildZkError(LogTag, "Error while collecting cpuUsageData")
		return response, respErr
	}
	cpuUsage := promResponse.ConvertMetricToPodUsage(cpuUsageData)

	memUsageData, err := metricServerRepo.GetPodMemoryUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		respErr := utils.BuildZkError(LogTag, "Error while collecting cpuUsageData")
		return response, respErr
	}
	memUsage := promResponse.ConvertMetricToPodUsage(memUsageData)

	response.CPUUsage = cpuUsage
	response.MemUsage = memUsage

	return response, nil
}

func (s prometheusService) TestIntegrationConnection(integrationId string) (promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		zkLogger.Error(LogTag, "Integration not found: ", integrationId, zkError)
		var resp promResponse.TestConnectionResponse
		resp.ConnectionStatus = utils.StatusError
		resp.ConnectionMessage = "Integration Not found"
		return resp, zkError
	}

	username, password := getUsernamePassword(*integration)
	resp, zkError := getConnectionStatus(integration.Url, username, password)
	if zkError != nil {
		return resp, zkError
	}

	metricServerResp, zkError := isIntegrationMetricServer(integrationId, integration.Url, username, password)
	if zkError != nil {
		return resp, zkError
	}

	resp.HasMetricServer = metricServerResp.MetricServer
	return resp, nil
}

func (s prometheusService) TestUnsavedIntegrationConnection(url string, username, password *string) (promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	resp, zkError := getConnectionStatus(url, username, password)
	if zkError != nil {
		return resp, zkError
	}

	metricServerResp, zkError := isIntegrationMetricServer("", url, username, password)
	if zkError != nil {
		return resp, zkError
	}

	resp.HasMetricServer = metricServerResp.MetricServer
	return resp, nil
}

func isIntegrationMetricServer(integrationId, url string, username, password *string) (promResponse.IsIntegrationMetricServerResponse, *zkerrors.ZkError) {
	var response promResponse.IsIntegrationMetricServerResponse
	resp, zkErr := getPrometheusApiResponse(url, username, password, utils.PrometheusQueryLabelValuesEndpoint, nil)
	if zkErr != nil {
		return response, zkErr
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId)
		response.StatusCode = common.ToPtr(resp.StatusCode)
		response.Status = common.ToPtr(resp.Status)
		response.Error = common.ToPtr(true)
		return response, nil
	}

	respBody, zkErr := getResponseBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	labelResponse, zkErr := readResponseBody[promResponse.StringListPrometheusResponse](respBody)

	for _, label := range labelResponse.Data {
		if strings.HasPrefix(label, "kubelet_") {
			response.MetricServer = common.ToPtr(true)
			return response, nil
		}
	}

	return response, nil
}

func (s prometheusService) GetMetricAttributes(integrationId string, metricName string, st string, et string) (promResponse.MetricAttributesListResponse, *zkerrors.ZkError) {
	var response promResponse.MetricAttributesListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		zkLogger.Error(LogTag, "Integration not found: ", integrationId, zkError)
		return response, zkError
	}

	username, password := getUsernamePassword(*integration)
	queryParam := map[string]string{
		"start":   st,
		"end":     et,
		"match[]": metricName,
	}

	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, utils.PrometheusQuerySeriesEndpoint, queryParam)
	if zkErr != nil {
		return response, zkErr
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId, resp.StatusCode, resp.Status)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return response, &zkErr
	}

	respBody, zkErr := getResponseBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	attributesResponse, zkErr := readResponseBody[promResponse.MetricAttributesPrometheusResponse](respBody)
	uniqueValueListPerAttribute := getUniqueValuesOfAttributes(attributesResponse.Data)
	response.Attributes = make(map[string]int)
	for key, value := range uniqueValueListPerAttribute {
		response.Attributes[key] = len(value)
	}

	return response, nil
}

func (s prometheusService) MetricsList(integrationId string) (promResponse.IntegrationMetricsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationMetricsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	username, password := getUsernamePassword(*integration)
	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, utils.PrometheusQueryLabelValuesEndpoint, nil)
	if zkErr != nil {
		return response, nil
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId, resp.StatusCode, resp.Status)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return response, &zkErr
	}

	respBody, zkErr := getResponseBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	metricsListResponse, zkErr := readResponseBody[promResponse.StringListPrometheusResponse](respBody)

	response.Metrics = metricsListResponse.Data
	return response, nil
}

func (s prometheusService) AlertsList(integrationId string, alertName string) (promResponse.IntegrationAlertsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationAlertsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	username, password := getUsernamePassword(*integration)
	if zkUtils.IsEmpty(alertName) {
		alertName = "!="
	}

	rulesResponse, zkErr := getAlertQuery(alertName, integration.Url, username, password)
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting alert query: ", zkErr)
		return response, zkErr
	}

	response = promResponse.ConvertAlertListPrometheusResponseToIntegrationAlertsListResponse(rulesResponse)
	return response, nil
}

func (s prometheusService) GetAlertsTimeSeriesTrigger(integrationId string, step string, startTime string, endTime string) (promResponse.AlertTimeSeriesResponse, *zkerrors.ZkError) {
	var response promResponse.AlertTimeSeriesResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	//http://localhost:9090/api/v1/query_range?query=ALERTS{alertstate='firing'}+OR+{alertstate='pending'}&start=1701004465.365&end=1701177265.365&step=691
	queryParam := fmt.Sprintf("query=ALERTS{alertstate='firing'}+OR+{alertstate='pending'}&start=%s&end=%s&step=%s", startTime, endTime, step)
	username, password := getUsernamePassword(*integration)
	//queryParam := map[string]string{
	//	"start": startTime,
	//	"end":   endTime,
	//	"step":  step,
	//	"query": "ALERTS{alertstate='firing'}+OR+{alertstate='pending'}",
	//}

	url := utils.PrometheusQueryAlertsRangeEndpoint + "?" + queryParam
	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, url, nil)
	if zkErr != nil {
		return response, zkErr
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId, resp.StatusCode, resp.Status)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return response, &zkErr
	}

	respBody, zkErr := getResponseBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	alertsResponse, zkErr := readResponseBody[promResponse.AlertRangePrometheusResponse](respBody)
	if zkErr != nil {
		return response, zkErr
	}

	stepInt, _ := strconv.ParseInt(step, 10, 64)

	alertNameToStateToRange := make(map[string]map[string][]promResponse.Duration)
	for _, alert := range alertsResponse.Data.Result {
		if alertNameToStateToRange[alert.Metric.AlertName] == nil {
			alertNameToStateToRange[alert.Metric.AlertName] = make(map[string][]promResponse.Duration)
		}
		alertNameToStateToRange[alert.Metric.AlertName][alert.Metric.AlertState] = findSeries(alert.Values, int(stepInt))
	}
	//alertsResponse = findSeries(alertsResponse.Data.Result, step)

	for alertName, stateToRange := range alertNameToStateToRange {
		alertRangeData := promResponse.AlertsRangeData{}
		alertRangeData.AlertName = alertName
		for state, rangeList := range stateToRange {
			seriesData := promResponse.SeriesData{}
			seriesData.State = state
			seriesData.Duration = rangeList
			alertRangeData.SeriesData = append(alertRangeData.SeriesData, seriesData)
		}
		response.AlertsRangeData = append(response.AlertsRangeData, alertRangeData)
	}

	return response, nil
}

func (s prometheusService) PrometheusAlertWebhook(integrationId string, alertWebhookData promResponse.AlertWebhookResponse) {
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		zkLogger.Error(LogTag, "Integration not found: ", integrationId, zkError)
	}

	username, password := getUsernamePassword(*integration)

	alertName := alertWebhookData.GroupLabels["alertname"].(string)
	rulesResponse, zkErr := getAlertQuery(alertName, integration.Url, username, password)
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting alert query: ", zkErr)
		return
	}
	var queryValues []string
	for _, group := range rulesResponse.Data.Groups {
		for _, rule := range group.Rules {
			queryValues = append(queryValues, rule.Query)
		}
	}

	var query string
	if len(queryValues) > 0 {
		query = queryValues[0]
	}

	zkLogger.Info(LogTag, "query: ", query)

	for i := range alertWebhookData.Alerts {
		alertWebhookData.Alerts[i].Query = query
	}

	jsonBody, err := json.Marshal(alertWebhookData)
	if err != nil {
		zkLogger.Error(LogTag, "Cannot Marshal data, encountered Err", err)
		return
	}

	bodyReader := bytes.NewReader(jsonBody)

	zkLogger.Info(LogTag, "request body is", string(jsonBody))

	response, zkErr := zkHttp.Create().Go("POST", "http://zk-gpt.zk-client.svc.cluster.local:80/v1/i/gpt/processPromAlert", bodyReader)
	if zkErr != nil {
		zkLogger.Error(LogTag, "error in making call to gpt", zkErr)
		return
	}

	statusCode := response.StatusCode
	if statusCode != 200 {
		status := response.Status
		zkLogger.Error(LogTag, "Status code not 200: ", statusCode, status)
		return
	}
	zkLogger.Info(LogTag, "Alert webhook sent successfully")
}

func getConnectionStatus(url string, username, password *string) (promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	var resp promResponse.TestConnectionResponse
	resp.ConnectionStatus = utils.StatusError

	queryParam := map[string]string{
		"query": "up",
	}

	httpResp, zkErr := getPrometheusApiResponse(url, username, password, utils.PrometheusQueryEndpoint, queryParam)
	if zkErr != nil {
		zkErrMetadata := zkErr.Metadata.(*zkerrors.ZkError)
		resp.ConnectionMessage = zkErrMetadata.Metadata.(string)
		return resp, nil
	}

	if httpResp.StatusCode != 200 {
		zkLogger.Info(LogTag, "Status code not 200")
		resp.ConnectionStatus = utils.StatusError
		resp.ConnectionMessage = httpResp.Status
		return resp, nil
	}

	respBody, zkErr := getResponseBody(httpResp)
	if zkErr != nil {
		resp.ConnectionMessage = "internal server error"
		return resp, zkErr
	}

	apiResponse, zkErr := readResponseBody[promResponse.QueryResultPrometheusResponse](respBody)
	if zkErr != nil {
		resp.ConnectionMessage = "internal server error"
		return resp, zkErr
	}

	if apiResponse.Status == utils.StatusSuccess {
		resp.ConnectionStatus = utils.StatusSuccess
		resp.ConnectionMessage = utils.ConnectionSuccessful
		return resp, nil
	} else if apiResponse.Status == utils.StatusError {
		resp.ConnectionMessage = apiResponse.Error
		return resp, nil
	}

	zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	return resp, &zkError
}

func getAlertQuery(alertName string, url string, username *string, password *string) (promResponse.AlertListPrometheusResponse, *zkerrors.ZkError) {
	//http://localhost:9090/api/v1/rules?type=alert&rule_name[]=InstanceDown
	var response promResponse.AlertListPrometheusResponse
	queryParam := map[string]string{
		"type":        "alert",
		"rule_name[]": alertName,
	}

	resp, zkErr := getPrometheusApiResponse(url, username, password, utils.PrometheusQueryRulesEndpoint, queryParam)
	if zkErr != nil {
		return response, zkErr
	}

	respBody, zkErr := getResponseBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	rulesResponse, zkErr := readResponseBody[promResponse.AlertListPrometheusResponse](respBody)
	if zkErr != nil {
		return response, zkErr
	}

	if rulesResponse.Status != utils.StatusSuccess {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return response, &zkErr
	}

	return rulesResponse, nil
}

func getIntegrationDetails(s prometheusService, integrationId string) (*dto.Integration, *zkerrors.ZkError) {
	var zkError *zkerrors.ZkError
	integration := s.integrationsManager.GetIntegrationById(integrationId)
	if integration == nil || integration.Id == "" {
		zkError = common.ToPtr(zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	return integration, zkError
}

func getUsernamePassword(integration dto.Integration) (*string, *string) {
	return integration.Authentication.Username, integration.Authentication.Password
}

func getPrometheusApiResponse(url string, username *string, password *string, prometheusQueryPath string, queryParams map[string]string) (*http.Response, *zkerrors.ZkError) {
	if common.IsEmpty(url) {
		zkLogger.Error(LogTag, "url is empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestEmptyUrl, nil)
		return nil, &zkError
	}

	//url = "http://localhost:9090"
	zkLogger.Info(LogTag, "url: ", url+prometheusQueryPath)

	req := zkHttp.Create()
	if queryParams != nil {
		for key, value := range queryParams {
			req = req.QueryParam(key, value)
		}
	}

	httpResp, zkErr := req.
		BasicAuth(username, password).
		Get(url + prometheusQueryPath)

	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while calling the api: ", url+prometheusQueryPath)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		return nil, &newZkErr
	}

	return httpResp, nil
}

func getUniqueValuesOfAttributes(attributes []promResponse.AttributesMap) map[string]ds.Set[string] {
	uniqueValues := make(map[string]ds.Set[string])
	for _, attribute := range attributes {
		for key, value := range attribute {
			if uniqueValues[key] == nil {
				uniqueValues[key] = ds.Set[string]{}
			}
			uniqueValues[key].Add(value)
		}
	}
	return uniqueValues
}

func getResponseBody(response *http.Response) ([]byte, *zkerrors.ZkError) {
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while reading the response body: ", err)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return nil, &zkErr
	}
	return respBody, nil
}

func readResponseBody[T any](responseBodyBytes []byte) (T, *zkerrors.ZkError) {
	var responseBody T
	err := json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return responseBody, &zkErr
	}

	return responseBody, nil
}

func extractMetricAttributes(dataVector model.Vector) promResponse.VectorList {
	var vectorList promResponse.VectorList = make([]promResponse.AttributesMap, 0)
	for _, sample := range dataVector {
		var attributes = make(map[string]string)
		for key, value := range sample.Metric {
			if strings.HasPrefix(string(key), "__") {
				continue
			}
			attributes[string(key)] = string(value)
		}
		vectorList = append(vectorList, attributes)
	}
	return vectorList
}

func findSeries(arr []promResponse.Value, step int) []promResponse.Duration {
	var result []promResponse.Duration

	if len(arr) == 0 {
		return result
	}

	// Sort the array based on the first element of each sub-array
	sort.Slice(arr, func(i, j int) bool {
		return int(arr[i][0].(float64)) < int(arr[j][0].(float64))
	})

	var start, end int

	for i := 0; i < len(arr); i++ {
		timestamp := int(arr[i][0].(float64))
		if i == 0 {
			start = timestamp
			end = timestamp
		} else {
			if timestamp-end == step {
				end = timestamp
			} else {
				result = append(result, promResponse.Duration{From: start, To: end})
				start = timestamp
				end = timestamp
			}
		}
	}

	result = append(result, promResponse.Duration{From: start, To: end})

	return result
}
