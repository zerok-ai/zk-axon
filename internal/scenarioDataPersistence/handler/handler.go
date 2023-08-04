package handler

import (
	"axon/internal/config"
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/service"
	"axon/internal/scenarioDataPersistence/validation"
	"axon/utils"
	zkErrorsScenarioManager "axon/utils/zkerrors"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

type TracePersistenceHandler interface {
	GetIssuesListWithDetailsHandler(ctx iris.Context)
	GetScenarioDetailsHandler(ctx iris.Context)
	GetIssueDetailsHandler(ctx iris.Context)
	GetIncidentListHandler(ctx iris.Context)
	GetIncidentDetailsHandler(ctx iris.Context)
	GetSpanRawDataHandler(ctx iris.Context)
	GetIncidentListForScenarioId(ctx iris.Context)
}

var LogTag = "trace_persistence_handler"

type tracePersistenceHandler struct {
	service service.TracePersistenceService
	cfg     config.AppConfigs
}

func (t tracePersistenceHandler) GetIncidentListForScenarioId(ctx iris.Context) {
	scenarioId := ctx.Params().Get(utils.ScenarioId)
	issueHash := ctx.URLParam(utils.IssueHashQueryParam)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.ValidateScenarioIdOffsetAndLimit(scenarioId, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentListForScenarioId api", err)
		zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IncidentIdListResponse](200, traceResponse.IncidentIdListResponse{}, nil, err)
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIncidentListServiceForScenarioId(scenarioId, issueHash, o, l)

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IncidentDetailListResponse]

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailListResponse](200, resp, resp, err)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailListResponse](200, resp, nil, err)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func NewTracePersistenceHandler(persistenceService service.TracePersistenceService, cfg config.AppConfigs) TracePersistenceHandler {
	return &tracePersistenceHandler{
		service: persistenceService,
		cfg:     cfg,
	}
}

func (t tracePersistenceHandler) GetIssuesListWithDetailsHandler(ctx iris.Context) {
	services := ctx.URLParam(utils.ServicesQueryParam)
	scenarioIds := ctx.URLParam(utils.ScenarioIdListQueryParam)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")
	st := ctx.URLParam(utils.StartTimeQueryParam)

	if err := validation.GetIssuesListWithDetails(offset, limit, st); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIssuesListWithDetailsHandler: ", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIssueListWithDetailsService(services, scenarioIds, st, l, o)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IssueListWithDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetScenarioDetailsHandler(ctx iris.Context) {
	services := ctx.URLParam(utils.ServicesQueryParam)
	scenarioIds := ctx.URLParam(utils.ScenarioIdListQueryParam)
	st := ctx.URLParam(utils.StartTimeQueryParam)

	if err := validation.ValidateGetScenarioDetails(scenarioIds, st); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetScenarioDetailsHandler: ", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, err := t.service.GetScenarioDetailsService(scenarioIds, services, st)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.ScenarioDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIssueDetailsHandler(ctx iris.Context) {
	issueHash := ctx.Params().Get(utils.IssueHash)

	if zkCommon.IsEmpty(issueHash) {
		zkLogger.Error(LogTag, "IssueHash is empty in GetIssueDetailsHandler api")
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestIssueHashEmpty, nil)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(zkErr.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, err := t.service.GetIssueDetailsService(issueHash)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IssueDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIncidentListHandler(ctx iris.Context) {
	issueHash := ctx.Params().Get(utils.IssueHash)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.ValidateIssueHashOffsetAndLimit(issueHash, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentListHandler api", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIncidentListService(issueHash, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IncidentIdListResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIncidentDetailsHandler(ctx iris.Context) {
	traceId := ctx.Params().Get(utils.IncidentId)
	spanId := ctx.URLParam(utils.SpanIdQueryParam)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.ValidateGetIncidentDetailsApi(traceId, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentDetailsHandler api", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIncidentDetailsService(traceId, spanId, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IncidentDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetSpanRawDataHandler(ctx iris.Context) {
	traceId := ctx.Params().Get(utils.IncidentId)
	spanId := ctx.Params().Get(utils.SpanId)
	if err := validation.ValidateGetSpanRawDataApi(traceId, spanId); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetSpanRawDataHandler api", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, err := t.service.GetSpanRawDataService(traceId, spanId)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.SpanRawDataResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
