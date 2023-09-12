package handler

import (
	"axon/internal/config"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	prometheusService "axon/internal/prometheus/service"
	tracePersistenceService "axon/internal/scenarioDataPersistence/service"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"strings"
	"time"
)

type PrometheusHandler interface {
	GetPodsInfoHandler(ctx iris.Context)
	GetContainerInfoHandler(ctx iris.Context)
	GetContainerMetricsHandler(ctx iris.Context)
	GetGenericQueryHandler(ctx iris.Context)
}

var LogTag = "prometheus_handler"

type prometheusHandler struct {
	tracePersistenceSvc tracePersistenceService.TracePersistenceService
	prometheusSvc       prometheusService.PrometheusService
	cfg                 config.AppConfigs
}

func NewPrometheusHandler(persistenceService prometheusService.PrometheusService,
	tracePersistenceSvc tracePersistenceService.TracePersistenceService,
	cfg config.AppConfigs) PrometheusHandler {
	return &prometheusHandler{
		tracePersistenceSvc: tracePersistenceSvc,
		prometheusSvc:       persistenceService,
		cfg:                 cfg,
	}
}

func (t prometheusHandler) GetGenericQueryHandler(ctx iris.Context) {
	var req request.GenericHTTPRequest
	readError := ctx.ReadJSON(&req)
	if readError != nil {
		zkLogger.Error(LogTag, "Error while reading request body: ", readError)
		ctx.StatusCode(500)
		return
	}
	promQuery := req.Query
	promDatasourceId := req.PromDatasourceId
	startTime := req.Time
	if startTime == 0 {
		startTime = time.Now().Unix()
	}
	durationStr := ctx.URLParamDefault(utils.DurationQueryParam, "0m")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing duration: ", err)
		ctx.StatusCode(500)
		return
	}

	endTime := startTime + int64(duration.Seconds())
	queryReq := request.GenericRequest{
		PromDatasourceId: promDatasourceId,
		Query:            string(promQuery),
		StartTime:        int64(startTime),
		EndTime:          int64(endTime),
		Duration:         int64(duration),
	}
	resp, zkErr := t.prometheusSvc.GetGenericQueryService(queryReq)

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.GenericQueryResponse]
	sendResponse[promResponse.GenericQueryResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetPodsInfoHandler(ctx iris.Context) {
	spansList, err := t.tracePersistenceSvc.GetIncidentDetailsService(ctx.Params().Get(utils.TraceId), "", 0, 50)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting spanList: ", err)
		ctx.StatusCode(500)
		return
	}
	podsMap := make(map[string]bool)
	spamItems := spansList.Spans
	for _, spanItems := range spamItems {
		if spanItems.Source != "" {
			podsMap[spanItems.Source] = true
		}
		if spanItems.Destination == "" {
			podsMap[spanItems.Destination] = true
		}
	}
	podsList := []string{}
	nsList := []string{}
	for k := range podsMap {
		podNameParts := strings.Split(k, "/")
		if len(podNameParts) != 2 {
			continue
		}
		namespace := podNameParts[0]
		podName := podNameParts[1]
		podsList = append(podsList, podName+".*")
		nsList = append(nsList, namespace)
	}
	zkLogger.Debug(LogTag, "podsList: ", podsList)
	podInfoReq := generatePromRequest(ctx)
	podInfoReq.Timestamp = time.Now().Unix()
	podInfoReq.PodsListStr = arrayToPromList(podsList)
	podInfoReq.NamespaceListStr = arrayToPromList(nsList)
	zkLogger.Debug(LogTag, "podsListStr: ", podInfoReq.PodsListStr)
	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.PodsInfoResponse]
	resp, zkErr := t.prometheusSvc.GetPodsInfoService(podInfoReq)
	sendResponse[promResponse.PodsInfoResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetContainerInfoHandler(ctx iris.Context) {
	podInfoReq := generatePromRequest(ctx)
	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.ContainerInfoResponse]
	resp, zkErr := t.prometheusSvc.GetContainerInfoService(podInfoReq)
	sendResponse[promResponse.ContainerInfoResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetContainerMetricsHandler(ctx iris.Context) {
	podInfoReq := generatePromRequest(ctx)
	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.ContainerMetricsResponse]
	resp, zkErr := t.prometheusSvc.GetContainerMetricService(podInfoReq)
	sendResponse[promResponse.ContainerMetricsResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}
