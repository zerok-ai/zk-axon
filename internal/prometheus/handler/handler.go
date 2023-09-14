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
	queryReq := generateGenericPromRequest(ctx, req)
	resp, zkErr := t.prometheusSvc.GetGenericQueryService(*queryReq)

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.GenericQueryResponse]
	sendResponse[promResponse.GenericQueryResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetPodsInfoHandler(ctx iris.Context) {
	traceId := ctx.Params().Get(utils.TraceId)
	podsList, nsList, err := getPodsAndNSListFromTrace(traceId, t.tracePersistenceSvc)
	if err != nil {
		ctx.StatusCode(500)
		return
	}
	zkLogger.Debug(LogTag, "podsList: ", podsList)

	podInfoReq := generatePromRequestMetadata(ctx)
	podInfoReq.Timestamp = time.Now().Unix()
	podInfoReq.PodsListStr = arrayToPromList(podsList)
	podInfoReq.NamespaceListStr = arrayToPromList(nsList)
	zkLogger.Debug(LogTag, "podsListStr: ", podInfoReq.PodsListStr)
	zkLogger.Debug(LogTag, "namespaceListStr: ", podInfoReq.NamespaceListStr)

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.PodsInfoResponse]
	resp, zkErr := t.prometheusSvc.GetPodsInfoService(podInfoReq)
	sendResponse[promResponse.PodsInfoResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetContainerInfoHandler(ctx iris.Context) {
	podInfoReq := generatePromRequestMetadata(ctx)
	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.ContainerInfoResponse]
	resp, zkErr := t.prometheusSvc.GetContainerInfoService(podInfoReq)
	sendResponse[promResponse.ContainerInfoResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}

func (t prometheusHandler) GetContainerMetricsHandler(ctx iris.Context) {
	podInfoReq := generatePromRequestMetadata(ctx)
	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.ContainerMetricsResponse]
	resp, zkErr := t.prometheusSvc.GetContainerMetricService(podInfoReq)
	sendResponse[promResponse.ContainerMetricsResponse](ctx, resp, zkHttpResponse, zkErr, t.cfg.Http.Debug)
}
