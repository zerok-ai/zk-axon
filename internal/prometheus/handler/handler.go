package handler

import (
	"axon/internal/config"
	"axon/internal/integrations/dto"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	prometheusService "axon/internal/prometheus/service"
	tracePersistenceService "axon/internal/scenarioDataPersistence/service"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"time"
)

type PrometheusHandler interface {
	GetPodsInfoHandler(ctx iris.Context)
	GetContainerInfoHandler(ctx iris.Context)
	GetContainerMetricsHandler(ctx iris.Context)
	GetGenericQueryHandler(ctx iris.Context)
	TestIntegrationConnectionStatus(ctx iris.Context)
	TestUnsavedIntegrationConnectionStatus(ctx iris.Context)
	GetMetrics(ctx iris.Context)
	GetMetricAttributes(ctx iris.Context)
	GetAlerts(ctx iris.Context)
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

func (t prometheusHandler) TestIntegrationConnectionStatus(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.TestConnectionResponse]
	var zkErr *zkerrors.ZkError
	resp, zkErr := t.prometheusSvc.TestIntegrationConnection(integrationId)

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.TestConnectionResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.TestConnectionResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t prometheusHandler) TestUnsavedIntegrationConnectionStatus(ctx iris.Context) {
	var req dto.UnsavedIntegrationRequestBody
	readError := ctx.ReadJSON(&req)
	if readError != nil {
		zkLogger.Error(LogTag, "Error while reading request body: ", readError)
		ctx.StatusCode(500)
		return
	}

	url := req.Url
	username := req.Username
	password := req.Password

	if zkCommon.IsEmpty(url) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Url is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.TestConnectionResponse]
	var zkErr *zkerrors.ZkError
	resp, zkErr := t.prometheusSvc.TestUnsavedIntegrationConnection(url, username, password)

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.TestConnectionResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.TestConnectionResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)

}

func (t prometheusHandler) GetMetricAttributes(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}
	metricName := ctx.Params().Get(utils.MetricAttributeNamePathParam)
	if zkCommon.IsEmpty(metricName) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("MetricName is required")
		return
	}

	// if start and end time is not passed, we will take data from last 6 hrs
	startTime := ctx.URLParam(utils.StartTimeQueryParam)
	endTime := ctx.URLParam(utils.EndTimeQueryParam)
	if zkCommon.IsEmpty(startTime) || zkCommon.IsEmpty(endTime) {
		currentTime := time.Now()
		startTime = strconv.FormatInt(currentTime.Unix(), 10)
		endTime = strconv.FormatInt(currentTime.Add(-(6 * time.Hour)).Unix(), 10)
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.MetricAttributesListResponse]
	var zkErr *zkerrors.ZkError
	resp, zkErr := t.prometheusSvc.GetMetricAttributes(integrationId, metricName, startTime, endTime)

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.MetricAttributesListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.MetricAttributesListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t prometheusHandler) GetMetrics(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.IntegrationMetricsListResponse]
	var zkErr *zkerrors.ZkError
	resp, zkErr := t.prometheusSvc.MetricsList(integrationId)

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.IntegrationMetricsListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.IntegrationMetricsListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t prometheusHandler) GetAlerts(ctx iris.Context) {
	integrationId := ctx.Params().Get(utils.IntegrationIdxPathParam)
	if zkCommon.IsEmpty(integrationId) {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("IntegrationId is required")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.IntegrationAlertsListResponse]
	var zkErr *zkerrors.ZkError
	resp, zkErr := t.prometheusSvc.AlertsList(integrationId)

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.IntegrationAlertsListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.IntegrationAlertsListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
