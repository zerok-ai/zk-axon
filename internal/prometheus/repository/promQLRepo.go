package repository

import (
	"axon/internal/prometheus/model/request"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	logger "github.com/zerok-ai/zk-utils-go/logs"
)

var LogTag = "zk_promQL_repo"

type PromQLRepo interface {
	GetPodCPUUsage(podInfoReq request.PromRequestMeta) (model.Matrix, error)
	GetPodMemoryUsage(podInfoReq request.PromRequestMeta) (model.Matrix, error)
	PodsInfoQuery(podInfoReq request.PromRequestMeta) (model.Vector, error)
	PodCreatedQuery(podInfoReq request.PromRequestMeta) (model.Vector, error)
	PodContainerInfoQuery(podInfoReq request.PromRequestMeta) (model.Vector, error)
}

const (
	CPUUsageQueryTemplate    = `sum(rate(container_cpu_usage_seconds_total{namespace="{{.Namespace}}", pod=~"{{.Pod}}", image!="", container!=""}[{{.RateInterval}}])) by (container)`
	MemoryUsageQueryTemplate = `sum(container_memory_working_set_bytes{namespace="{{.Namespace}}", pod=~"{{.Pod}}", image!="", container!=""}) by (container)`
	PodsInfoQuery            = `kube_pod_info{pod=~"^({{.PodsListStr}})$"} @ {{.Timestamp}}`
	PodCreatedQuery          = `kube_pod_created{namespace="{{.Namespace}}",pod=~"{{.Pod}}"} @ {{.Timestamp}}`
	PodContainerInfoQuery    = `kube_pod_container_info{namespace="{{.Namespace}}",pod=~"{{.Pod}}"} @ {{.Timestamp}}`
)

type promQLRepo struct {
	promClient api.Client
	queryAPI   v1.API
}

func NewPromQLRepo(client api.Client) PromQLRepo {
	return &promQLRepo{
		promClient: client,
		queryAPI:   v1.NewAPI(client),
	}
}

func (r promQLRepo) PodsInfoQuery(podInfoReq request.PromRequestMeta) (model.Vector, error) {
	query, err := GetPromQueryString(PodsInfoQuery, podInfoReq)
	if err != nil {
		return nil, err
	}
	logger.Debug(LogTag, "Query: ", query)
	podInfo, err := r.GetPromData(query, podInfoReq.EndTime)
	if err != nil {
		return nil, err
	}
	return podInfo, nil
}

func (r promQLRepo) PodCreatedQuery(podInfoReq request.PromRequestMeta) (model.Vector, error) {
	query, err := GetPromQueryString(PodCreatedQuery, podInfoReq)
	if err != nil {
		return nil, err
	}
	logger.Debug(LogTag, "Query: ", query)
	podCreated, err := r.GetPromData(query, podInfoReq.EndTime)
	if err != nil {
		return nil, err
	}
	return podCreated, nil
}

func (r promQLRepo) PodContainerInfoQuery(podInfoReq request.PromRequestMeta) (model.Vector, error) {
	query, err := GetPromQueryString(PodContainerInfoQuery, podInfoReq)
	if err != nil {
		return nil, err
	}
	logger.Debug(LogTag, "Query: ", query)
	podContainerInfo, err := r.GetPromData(query, podInfoReq.EndTime)
	if err != nil {
		return nil, err
	}
	return podContainerInfo, nil
}

func (r promQLRepo) GetPodCPUUsage(podInfoReq request.PromRequestMeta) (model.Matrix, error) {
	query, err := GetPromQueryString(CPUUsageQueryTemplate, podInfoReq)
	if err != nil {
		return nil, err
	}
	logger.Debug(LogTag, "Query: ", query)
	logger.Debug(LogTag, podInfoReq.Namespace, podInfoReq.Pod)
	cpuMetric, err := r.GetPromMatrixData(query, podInfoReq.StartTime, podInfoReq.EndTime)
	if err != nil {
		return nil, err
	}
	return cpuMetric, nil
}

func (r promQLRepo) GetPodMemoryUsage(podInfoReq request.PromRequestMeta) (model.Matrix, error) {
	query, err := GetPromQueryString(MemoryUsageQueryTemplate, podInfoReq)
	if err != nil {
		return nil, err
	}
	logger.Debug(LogTag, "Query: ", query)
	memoryMetric, err := r.GetPromMatrixData(query, podInfoReq.StartTime, podInfoReq.EndTime)
	if err != nil {
		return nil, err
	}
	return memoryMetric, nil
}
