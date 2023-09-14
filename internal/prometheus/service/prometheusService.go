package service

import (
	"axon/internal/integrations"
	"axon/internal/integrations/dto"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/repository"
	zkUtils "axon/utils"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/model"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
)

var LogTag = "zk_prometheus_service"

type PrometheusService interface {
	GetPodsInfoService(podInfoReq request.PromRequestMeta) (promResponse.PodsInfoResponse, *zkErrors.ZkError)
	GetContainerInfoService(podInfoReq request.PromRequestMeta) (promResponse.ContainerInfoResponse, *zkErrors.ZkError)
	GetContainerMetricService(podInfoReq request.PromRequestMeta) (promResponse.ContainerMetricsResponse, *zkErrors.ZkError)
	GetGenericQueryService(genericQueryReq request.GenericPromRequest) (promResponse.GenericQueryResponse, *zkErrors.ZkError)

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

func (s prometheusService) GetGenericQueryService(genericQueryReq request.GenericPromRequest) (promResponse.GenericQueryResponse, *zkErrors.ZkError) {
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

func (s prometheusService) GetPodsInfoService(podInfoReq request.PromRequestMeta) (promResponse.PodsInfoResponse, *zkErrors.ZkError) {
	var response promResponse.PodsInfoResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podsInfo, err := metricServerRepo.PodsInfoQuery(podInfoReq)
	if err != nil {
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting podInfo: ", err.Error())
		return response, respErr
	}

	podsInfoItems := extractMetricAttributes(podsInfo)
	response.PodsInfo = podsInfoItems
	return response, nil
}

func (s prometheusService) GetContainerInfoService(podInfoReq request.PromRequestMeta) (promResponse.ContainerInfoResponse, *zkErrors.ZkError) {
	var response promResponse.ContainerInfoResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}
	podContainerInfo, err := metricServerRepo.PodContainerInfoQuery(podInfoReq)
	if err != nil {
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting podContainerInfo: ", err.Error())
		return response, respErr
	}
	podContainerInfoItems := extractMetricAttributes(podContainerInfo)
	response.ContainerInfo = podContainerInfoItems
	return response, nil
}

func (s prometheusService) GetContainerMetricService(podInfoReq request.PromRequestMeta) (promResponse.ContainerMetricsResponse, *zkErrors.ZkError) {
	var response promResponse.ContainerMetricsResponse
	metricServerRepo := s.GetMetricServerRepo()
	if metricServerRepo == nil {
		respErr := zkUtils.BuildZkError(LogTag, "No metric server found")
		return response, respErr
	}

	cpuUsageData, err := metricServerRepo.GetPodCPUUsage(podInfoReq)
	if err != nil {
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		return response, respErr
	}
	cpuUsage := promResponse.ConvertMetricToPodUsage(cpuUsageData)

	memUsageData, err := metricServerRepo.GetPodMemoryUsage(podInfoReq)
	if err != nil {
		respErr := zkUtils.BuildZkError(LogTag, "Error while collecting cpuUsageData: ", err.Error())
		return response, respErr
	}
	memUsage := promResponse.ConvertMetricToPodUsage(memUsageData)

	response.CPUUsage = cpuUsage
	response.MemUsage = memUsage

	return response, nil
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
