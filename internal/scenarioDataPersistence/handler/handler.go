package handler

import (
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
	GetIssueDetailsHandler(ctx iris.Context)
	GetIncidentListHandler(ctx iris.Context)
	GetIncidentDetailsHandler(ctx iris.Context)
	GetSpanRawDataHandler(ctx iris.Context)
	//GetAllScenariosTracesData(ctx iris.Context)
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
	issueId := ctx.Params().Get(utils.IssueId)

	if zkCommon.IsEmpty(issueId) {
		zkLogger.Error(LogTag, "IssueHash is empty in GetIssueDetailsHandler api")
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestIssueIdEmpty, nil)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(zkErr.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	resp, err := t.service.GetIssueDetailsService(issueId)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IssueWithDetailsResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIncidentListHandler(ctx iris.Context) {
	issueId := ctx.Params().Get(utils.IssueId)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	if err := validation.ValidateGetIncidentApi(issueId, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentListHandler api", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIncidentListService(issueId, o, l)

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
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")
	if err := validation.ValidateGetSpanRawDataApi(traceId, spanId, offset, limit); err != nil {
		zkLogger.Error(LogTag, "Error while validating GetSpanRawDataHandler api", err)
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetSpanRawDataService(traceId, spanId, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.SpanRawDataResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

//func (t tracePersistenceHandler) GetAllScenariosTracesData(ctx iris.Context) {
//	scenarioId := ctx.Params().Get(utils.Scenario)
//	limit := ctx.URLParamDefault(utils.LimitQueryParam, "1000")
//	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")
//	if err := validation.ValidateGetScenariosAllTraceDataApi(scenarioId, offset, limit); err != nil {
//		zkLogger.Error(LogTag, "Error while validating GetAllScenariosTracesData api", err)
//		z := &zkHttp.ZkHttpResponseBuilder[any]{}
//		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
//		ctx.StatusCode(zkHttpResponse.Status)
//		ctx.JSON(zkHttpResponse)
//		return
//	}
//
//	l, _ := strconv.Atoi(limit)
//	o, _ := strconv.Atoi(offset)
//
//	resp, err := t.service.GetScenariosAllTracesDataService(scenarioId, o, l)
//
//	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.ScenarioIncidentDetailsResponse](200, resp, resp, err)
//	ctx.StatusCode(zkHttpResponse.Status)
//	ctx.JSON(zkHttpResponse)
//}
