package handler

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/service"
	"axon/internal/scenarioDataPersistence/validation"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"strconv"
)

type TracePersistenceHandler interface {
	GetIssuesListWithDetailsHandler(ctx iris.Context)
	GetIssueDetailsHandler(ctx iris.Context)
	GetIncidentListHandler(ctx iris.Context)
	GetIncidentDetailsHandler(ctx iris.Context)
	GetSpanRawDataHandler(ctx iris.Context)
}

var LogTag = "trace_persistence_handler"

type tracePersistenceHandler struct {
	service service.TracePersistenceService
}

func NewTracePersistenceHandler(persistenceService service.TracePersistenceService) TracePersistenceHandler {
	return &tracePersistenceHandler{
		service: persistenceService,
	}
}

func (t tracePersistenceHandler) GetIssuesListWithDetailsHandler(ctx iris.Context) {
	services := ctx.URLParam(utils.ServicesQueryParam)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.GetIssuesListWithDetails(offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIssuesListWithDetailsHandler: ", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIssueListWithDetailsService(services, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IssueListWithDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIssueDetailsHandler(ctx iris.Context) {
	issueHash := ctx.Params().Get(utils.IssueHash)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.ValidateIssueHashOffsetAndLimit(issueHash, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIssueDetailsHandler: ", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, err := t.service.GetIssueDetailsService(issueHash)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IssueListWithDetailsResponse](200, resp, resp, err)
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

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IncidentListResponse](200, resp, resp, err)
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
