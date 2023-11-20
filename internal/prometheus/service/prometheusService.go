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
	TestIntegrationConnection(integrationId string) (*promResponse.TestConnectionResponse, *zkerrors.ZkError)
	TestUnsavedIntegrationConnection(url, userName, password string) (*promResponse.TestConnectionResponse, *zkerrors.ZkError)
	IsIntegrationMetricServer(integrationId string) (promResponse.IsIntegrationMetricServerResponse, *zkerrors.ZkError)
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

func (s prometheusService) TestIntegrationConnection(integrationId string) (*promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	var resp promResponse.TestConnectionResponse
	resp.Status = "error"

	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		resp.Message = "Integration Not found"
		return &resp, zkError
	}

	username, password := getUsernamePassword(*integration)
	httpResp, zkErr := getPrometheusApiResponse(integration.Url, username, password, "/api/v1/query?query=up")
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkErr)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while reading the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	var apiResponse promResponse.QueryResult
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	if apiResponse.Status == "success" {
		resp.Status = "success"
		resp.Message = "Connection successful"
		return &resp, nil
	} else if apiResponse.Status == "error" {
		resp.Message = apiResponse.Error
		resp.ErrorType = apiResponse.ErrorType
		return &resp, nil
	}

	resp.Message = apiResponse.Error
	resp.Message = apiResponse.ErrorType
	return &resp, nil
}

func (s prometheusService) TestUnsavedIntegrationConnection(url, userName, password string) (*promResponse.TestConnectionResponse, *zkerrors.ZkError) {
	var resp promResponse.TestConnectionResponse
	resp.Status = "error"

	httpResp, zkErr := getPrometheusApiResponse(url, userName, password, "/api/v1/query?query=up")
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the integration status: ", zkErr)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while reading the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	var apiResponse promResponse.QueryResult
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		resp.Message = "internal server error"
		return &resp, &newZkErr
	}

	if apiResponse.Status == "success" {
		resp.Status = "success"
		resp.Message = "Connection successful"
		return &resp, nil
	} else if apiResponse.Status == "error" {
		resp.Message = apiResponse.Error
		resp.ErrorType = apiResponse.ErrorType
		return &resp, nil
	}

	resp.Message = apiResponse.Error
	resp.Message = apiResponse.ErrorType
	return &resp, nil

}

func (s prometheusService) IsIntegrationMetricServer(integrationId string) (promResponse.IsIntegrationMetricServerResponse, *zkerrors.ZkError) {
	var response promResponse.IsIntegrationMetricServerResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	resp, zkErr := getPrometheusResponse(*integration, "/api/v1/label/__name__/values")
	if zkErr != nil {
		return response, zkErr
	}

	var labelResponse promResponse.LabelNameResponse
	err := json.Unmarshal(resp, &labelResponse)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return response, &newZkErr
	}

	for _, label := range labelResponse.Data {
		if strings.HasPrefix(label, "kubelet_") {
			response.MetricServer = true
			return response, nil
		}
	}

	return response, nil
}

func (s prometheusService) MetricsList(integrationId string) (promResponse.IntegrationMetricsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationMetricsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	resp, zkErr := getPrometheusResponse(*integration, "/api/v1/label/__name__/values")
	if zkErr != nil {
		return response, zkErr
	}

	var metricsListResponse promResponse.LabelNameResponse
	err := json.Unmarshal(resp, &metricsListResponse)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return response, &newZkErr
	}

	response.Metrics = metricsListResponse.Data
	return response, nil
}

func (s prometheusService) AlertsList(integrationId string) (promResponse.IntegrationAlertsListResponse, *zkerrors.ZkError) {
	var response promResponse.IntegrationAlertsListResponse
	integration, zkError := getIntegrationDetails(s, integrationId)
	if zkError != nil {
		return response, zkError
	}

	resp, zkErr := getPrometheusResponse(*integration, "/api/v1/alerts")
	if zkErr != nil {
		return response, zkErr
	}

	var alertsResponse promResponse.AlertsResponse
	err := json.Unmarshal(resp, &alertsResponse)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return response, &newZkErr
	}

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

func getUsernamePassword(integration dto.Integration) (string, string) {
	return integration.Authentication.Username, integration.Authentication.Password
}

func getPrometheusResponse(integration dto.Integration, prometheusQueryPath string) ([]byte, *zkerrors.ZkError) {
	username, password := getUsernamePassword(integration)
	httpResp, zkErr := getPrometheusApiResponse(integration.Url, username, password, prometheusQueryPath)
	if zkErr != nil {
		zkLogger.Error(LogTag, "Error while getting the metrics status: ", zkErr)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, zkErr)
		return nil, &newZkErr
	}

	resp, err := io.ReadAll(httpResp.Body)
	if err != nil {
		zkLogger.Error(LogTag, "Error while reading the response body: ", err)
		newZkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, err)
		return nil, &newZkErr
	}

	return resp, nil
}

func getPrometheusApiResponse(url, username, password, prometheusQueryPath string) (*http.Response, *zkerrors.ZkError) {
	if common.IsEmpty(url) {
		zkLogger.Error(LogTag, "url is empty")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestEmptyUrl, nil)
		return nil, &zkError
	}

	zkLogger.Info(LogTag, "url: ", url)

	return zkHttp.Create().
		BasicAuth(username, password).
		Get(url + prometheusQueryPath)
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
