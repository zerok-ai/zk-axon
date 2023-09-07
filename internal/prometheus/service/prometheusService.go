package service

import (
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/repository"
	"github.com/prometheus/common/model"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

var LogTag = "zk_prometheus_service"

type PrometheusService interface {
	GetPodDetailsService(podInfoReq request.PodInfoRequest) (promResponse.PodDetailResponse, *zkErrors.ZkError)
}

func NewPrometheusService(repo repository.PromQLRepo) PrometheusService {
	return prometheusService{repo: repo}
}

type prometheusService struct {
	repo repository.PromQLRepo
}

func (s prometheusService) GetPodDetailsService(podInfoReq request.PodInfoRequest) (promResponse.PodDetailResponse, *zkErrors.ZkError) {
	var response promResponse.PodDetailResponse
	var cpuUsage promResponse.Usage
	var memUsage promResponse.Usage

	cpuUsageData, err := s.repo.GetPodCPUUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err)
		return response, nil
	}
	cpuUsage = promResponse.ConvertMetricToPodUsage("CPU Usage", cpuUsageData)

	memUsageData, err := s.repo.GetPodMemoryUsage(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting cpuUsageData: ", err)
		return response, nil
	}
	memUsage = promResponse.ConvertMetricToPodUsage("Memory Usage", memUsageData)

	podInfo, err := s.repo.PodInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podInfo: ", err)
		return response, nil
	}

	podCreated, err := s.repo.PodCreatedQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podCreated: ", err)
		return response, nil
	}

	podContainerInfo, err := s.repo.PodContainerInfoQuery(podInfoReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting podContainerInfo: ", err)
		return response, nil
	}

	podContainerInfoItems := extraceMetricAttributes(podContainerInfo)
	podCreatedItems := extraceMetricAttributes(podCreated)
	podInfoItems := extraceMetricAttributes(podInfo)

	var metadataItems = make(map[string]interface{})
	mergeMaps(metadataItems, podContainerInfoItems)
	mergeMaps(metadataItems, podCreatedItems)
	mergeMaps(metadataItems, podInfoItems)

	//podMetadata, err := s.repo.GetPodMetadata(namespace, podId)
	podMetadata := metadataItems

	response = promResponse.PodDetailResponse{
		PodName:      "podId",
		Metadata:     podMetadata,
		ZkInferences: "ZkInferences",
		CPUUsage:     cpuUsage,
		MemUsage:     memUsage,
	}

	return response, nil
}

func extraceMetricAttributes(dataVector model.Vector) map[string]interface{} {
	var attributes = make(map[string]interface{})
	for _, sample := range dataVector {
		for key, value := range sample.Metric {
			attributes[string(key)] = string(value)
		}
	}
	return attributes
}

func mergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}
