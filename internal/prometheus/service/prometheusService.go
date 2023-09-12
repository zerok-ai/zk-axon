package service

import (
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/repository"
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
	GetGenericQueryService(genericQueryReq request.GenericRequest) (promResponse.GenericQueryResponse, *zkErrors.ZkError)
}

func NewPrometheusService(metricServerRepo repository.PromQLRepo, dataSources map[string]repository.PromQLRepo) PrometheusService {
	return prometheusService{metricServerRepo: metricServerRepo, dataSources: dataSources}
}

type prometheusService struct {
	metricServerRepo repository.PromQLRepo
	dataSources      map[string]repository.PromQLRepo
}

func (s prometheusService) GetGenericQueryService(genericQueryReq request.GenericRequest) (promResponse.GenericQueryResponse, *zkErrors.ZkError) {
	var response promResponse.GenericQueryResponse
	queryResult, resultType, err := s.metricServerRepo.GenericQuery(genericQueryReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting queryResult: ", err)
		return response, nil
	}
	response.Result = queryResult
	response.Type = resultType
	return response, nil
}

func (s prometheusService) GetPodsInfoService(podInfoReq request.PromRequestMeta) (promResponse.PodsInfoResponse, *zkErrors.ZkError) {
	var response promResponse.PodsInfoResponse

	podsInfo, err := s.metricServerRepo.PodsInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podInfo: ", err)
		return response, nil
	}

	podsInfoItems := extractMetricAttributes(podsInfo)
	response.PodsInfo = podsInfoItems
	return response, nil
}

func (s prometheusService) GetContainerInfoService(podInfoReq request.PromRequestMeta) (promResponse.ContainerInfoResponse, *zkErrors.ZkError) {
	var response promResponse.ContainerInfoResponse
	podContainerInfo, err := s.metricServerRepo.PodContainerInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podContainerInfo: ", err)
		return response, nil
	}
	podContainerInfoItems := extractMetricAttributes(podContainerInfo)
	response.ContainerInfo = podContainerInfoItems
	return response, nil
}

func (s prometheusService) GetContainerMetricService(podInfoReq request.PromRequestMeta) (promResponse.ContainerMetricsResponse, *zkErrors.ZkError) {
	var response promResponse.ContainerMetricsResponse

	cpuUsageData, err := s.metricServerRepo.GetPodCPUUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err)
		return response, nil
	}
	cpuUsage := promResponse.ConvertMetricToPodUsage(cpuUsageData)

	memUsageData, err := s.metricServerRepo.GetPodMemoryUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err)
		return response, nil
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

func mergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}
