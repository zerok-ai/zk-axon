package repository

import (
	"axon/internal/prometheus/model/request"
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var LogTag = "zk_promQL_repo"

type PromQLRepo interface {
	GetPodCPUUsage(podInfoReq request.PodInfoRequest) (model.Matrix, error)
	GetPodMemoryUsage(podInfoReq request.PodInfoRequest) (model.Matrix, error)
	PodInfoQuery(podInfoReq request.PodInfoRequest) (model.Vector, error)
	PodCreatedQuery(podInfoReq request.PodInfoRequest) (model.Vector, error)
	PodContainerInfoQuery(podInfoReq request.PodInfoRequest) (model.Vector, error)
}

const (
	CPUUsageQueryTemplate    = `sum(rate(container_cpu_usage_seconds_total{namespace="{{.Namespace}}", pod=~"{{.Pod}}"}[{{.RateInterval}}])) by (container)`
	MemoryUsageQueryTemplate = `sum(container_memory_working_set_bytes{namespace="{{.Namespace}}", pod=~"{{.Pod}}"}) by (container)`
	PodInfoQuery             = `kube_pod_info{namespace="{{.Namespace}}",pod=~"{{.Pod}}"} @ {{.Timestamp}}`
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

func (r promQLRepo) PodInfoQuery(podInfoReq request.PodInfoRequest) (model.Vector, error) {
	query, err := r.GetPromQueryString(PodInfoQuery, podInfoReq)
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

func (r promQLRepo) PodCreatedQuery(podInfoReq request.PodInfoRequest) (model.Vector, error) {
	query, err := r.GetPromQueryString(PodCreatedQuery, podInfoReq)
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

func (r promQLRepo) PodContainerInfoQuery(podInfoReq request.PodInfoRequest) (model.Vector, error) {
	query, err := r.GetPromQueryString(PodContainerInfoQuery, podInfoReq)
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

func (r promQLRepo) GetPodCPUUsage(podInfoReq request.PodInfoRequest) (model.Matrix, error) {
	query, err := r.GetPromQueryString(CPUUsageQueryTemplate, podInfoReq)
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

func (r promQLRepo) GetPodMemoryUsage(podInfoReq request.PodInfoRequest) (model.Matrix, error) {
	query, err := r.GetPromQueryString(MemoryUsageQueryTemplate, podInfoReq)
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

func (r promQLRepo) GetPromQueryString(templateString string, podInfoReq request.PodInfoRequest) (string, error) {
	// Create a PromQL query
	queryTemplate, err := template.New("query").Parse(templateString)
	if err != nil {
		logger.Fatal(err.Error())
	}
	query := new(strings.Builder)
	err = queryTemplate.Execute(query, podInfoReq)

	// Query Prometheus
	logger.Debug(LogTag, "cpu query: ", query.String())
	logger.Debug(LogTag, "over: "+podInfoReq.StartTime.String()+" to "+podInfoReq.EndTime.String())
	return query.String(), nil
}

func (r promQLRepo) GetPromMatrixData(query string, startTime time.Time, endTime time.Time) (model.Matrix, error) {
	// Execute the query
	ctx := context.Background()
	result, warnings, err := r.queryAPI.QueryRange(ctx, query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  1 * time.Minute, // Adjust the step as needed
	})
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	logger.Debug(LogTag, "Result type: ", reflect.TypeOf(result).Name())

	// Process query result
	if matrix, ok := result.(model.Matrix); ok {
		return matrix, nil
	} else {
		logger.Debug(LogTag, "Query did not return a matrix\n")
	}

	return model.Matrix{}, nil
}

func (r promQLRepo) GetPromData(query string, endTime time.Time) (model.Vector, error) {
	// Execute the query
	ctx := context.Background()
	result, warnings, err := r.queryAPI.Query(ctx, query, endTime)
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	logger.Debug(LogTag, "Result type: ", reflect.TypeOf(result).Name())
	// Process query result
	if vector, ok := result.(model.Vector); ok {
		return vector, nil
	} else {
		logger.Debug(LogTag, "Query did not return a Vector\n")
	}

	return model.Vector{}, nil
}
