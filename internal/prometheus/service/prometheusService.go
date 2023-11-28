package service

import (
	"axon/internal/integrations"
	"axon/internal/integrations/dto"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/repository"
	zkUtils "axon/utils"
	zkErrorsAxon "axon/utils/zkerrors"
	"encoding/json"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/model"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/ds"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"net/http"
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
	AlertsList(integrationId string) (promResponse.IntegrationAlertsListResponse, *zkerrors.ZkError)

	GetMetricServerRepo() repository.PromQLRepo
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
		respErr := zkUtils.BuildZkError(LogTag, "No prometheus integration found for id: ", genericQueryReq.PromIntegrationId)
		return response, respErr
	}
	queryResult, resultType, err := targetPromIntegration.GenericQuery(genericQueryReq)
	if err != nil {
		respErr := zkUtils.BuildZkError(LogTag, "Failed to query prometheus, Error: ", err.Error())
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
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podsInfo, err := metricServerRepo.PodsInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podInfo: ", err.Error())
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting podInfo")
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
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podContainerInfo, err := metricServerRepo.PodContainerInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podContainerInfo: ", err.Error())
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting podContainerInfo")
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
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}

	cpuUsageData, err := metricServerRepo.GetPodCPUUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting cpuUsageData")
		return response, respErr
	}
	cpuUsage := promResponse.ConvertMetricToPodUsage(cpuUsageData)

	memUsageData, err := metricServerRepo.GetPodMemoryUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting cpuUsageData")
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
		resp.ConnectionStatus = zkUtils.StatusError
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

func getConnectionStatus(url string, username, password *string) (promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	var resp promResponse.TestConnectionResponse
	resp.ConnectionStatus = zkUtils.StatusError

	queryParam := map[string]string{
		"query": "up",
	}

	httpResp, zkErr := getPrometheusApiResponse(url, username, password, zkUtils.PrometheusQueryEndpoint, queryParam)
	if zkErr != nil {
		zkErrMetadata := zkErr.Metadata.(*zkerrors.ZkError)
		resp.ConnectionMessage = zkErrMetadata.Metadata.(string)
		return resp, nil
	}

	if httpResp.StatusCode != 200 {
		zkLogger.Info(LogTag, "Status code not 200")
		resp.ConnectionStatus = zkUtils.StatusError
		resp.ConnectionMessage = httpResp.Status
		return resp, nil
	}

	respBody, zkErr := getRequestBody(httpResp)
	if zkErr != nil {
		resp.ConnectionMessage = "internal server error"
		return resp, zkErr
	}

	apiResponse, zkErr := readResponseBody[promResponse.QueryResult](respBody)
	if zkErr != nil {
		resp.ConnectionMessage = "internal server error"
		return resp, zkErr
	}

	if apiResponse.Status == zkUtils.StatusSuccess {
		resp.ConnectionStatus = zkUtils.StatusSuccess
		resp.ConnectionMessage = zkUtils.ConnectionSuccessful
		return resp, nil
	} else if apiResponse.Status == zkUtils.StatusError {
		resp.ConnectionMessage = apiResponse.Error
		return resp, nil
	}

	zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
	return resp, &zkError
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
	resp, zkErr := getPrometheusApiResponse(url, username, password, zkUtils.PrometheusQueryLabelValuesEndpoint, nil)
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

	respBody, zkErr := getRequestBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	labelResponse, zkErr := readResponseBody[promResponse.LabelNameResponse](respBody)

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

	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, zkUtils.PrometheusQuerySeriesEndpoint, queryParam)
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

	respBody, zkErr := getRequestBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	attributesResponse, zkErr := readResponseBody[promResponse.MetricAttributes](respBody)
	uniqueValueListPerAttribute := getUniqueValuesOfAttributes(attributesResponse.Data)
	response.Attributes = make(map[string]int)
	for key, value := range uniqueValueListPerAttribute {
		response.Attributes[key] = len(value)
	}

	return response, nil
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

func (s prometheusService) MetricsList(integrationId string) (promResponse.IntegrationMetricsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationMetricsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	username, password := getUsernamePassword(*integration)
	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, zkUtils.PrometheusQueryLabelValuesEndpoint, nil)
	if zkErr != nil {
		return response, nil
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId)
		response.StatusCode = common.ToPtr(resp.StatusCode)
		response.Status = common.ToPtr(resp.Status)
		response.Error = common.ToPtr(true)
		response.Metrics = nil
		return response, nil
	}

	respBody, zkErr := getRequestBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	metricsListResponse, zkErr := readResponseBody[promResponse.LabelNameResponse](respBody)

	response.Metrics = metricsListResponse.Data
	return response, nil
}

func (s prometheusService) AlertsList(integrationId string) (promResponse.IntegrationAlertsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationAlertsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	username, password := getUsernamePassword(*integration)
	queryParam := map[string]string{
		"alertname":  "!=",
		"severity":   "critical",
		"alertstate": "firing",
	}

	resp, zkErr := getPrometheusApiResponse(integration.Url, username, password, zkUtils.PrometheusQueryAlertsEndpoint, queryParam)
	if zkErr != nil {
		return response, zkErr
	}

	if resp.StatusCode != 200 {
		zkLogger.Error(LogTag, "Status code not 200, integrationId: ", integrationId)
		response.StatusCode = common.ToPtr(resp.StatusCode)
		response.Status = common.ToPtr(resp.Status)
		response.Error = common.ToPtr(true)
		response.Alerts = nil
		return response, nil
	}

	respBody, zkErr := getRequestBody(resp)
	if zkErr != nil {
		return response, zkErr
	}

	alertsResponse, zkErr := readResponseBody[promResponse.AlertResponse](respBody)
	if zkErr != nil {
		return response, zkErr
	}

	alertList := make([]string, 0)
	for _, alert := range alertsResponse.Data.Alerts {
		alertList = append(alertList, alert.Labels["alertname"])
	}

	response.Alerts = alertList
	return response, nil
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

func getRequestBody(response *http.Response) ([]byte, *zkerrors.ZkError) {
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
