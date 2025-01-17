package handler

import (
	"axon/internal/config"
	"axon/internal/scenarioDataPersistence/model/request"
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/service"
	"axon/internal/scenarioDataPersistence/validation"
	"axon/utils"
	zkErrorsAxon "axon/utils/zkerrors"
	"encoding/json"
	"github.com/kataras/iris/v12"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrorsUtils "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

type TracePersistenceHandler interface {
	GetIssuesListWithDetailsHandler(ctx iris.Context)
	GetScenarioDetailsHandler(ctx iris.Context)
	GetIssueDetailsHandler(ctx iris.Context)
	GetIncidentListHandler(ctx iris.Context)
	GetPodDetailsHandler(ctx iris.Context)
	GetSpanRawDataHandler(ctx iris.Context)
	GetIncidentListForScenarioId(ctx iris.Context)
	GetErrorDataHandler(ctx iris.Context)
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

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IncidentDetailListResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.IncidentDetailListResponse

	if zkErr := validation.ValidateScenarioIdOffsetAndLimit(scenarioId, offset, limit); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentListForScenarioId api", zkErr)
	} else {
		l, _ := strconv.Atoi(limit)
		o, _ := strconv.Atoi(offset)
		resp, zkErr = t.service.GetIncidentListServiceForScenarioId(scenarioId, issueHash, o, l)
	}

	if zkErr == nil && resp.TotalRecords == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailListResponse](200, resp, nil, zkErr)
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

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IssueListWithDetailsResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.IssueListWithDetailsResponse

	if zkErr := validation.GetIssuesListWithDetails(offset, limit, st); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetIssuesListWithDetailsHandler: ", zkErr)
	} else {
		l, _ := strconv.Atoi(limit)
		o, _ := strconv.Atoi(offset)
		resp, zkErr = t.service.GetIssueListWithDetailsService(services, scenarioIds, st, l, o)
	}

	// DONE
	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IssueListWithDetailsResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IssueListWithDetailsResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetScenarioDetailsHandler(ctx iris.Context) {
	services := ctx.URLParam(utils.ServicesQueryParam)
	scenarioIds := ctx.URLParam(utils.ScenarioIdListQueryParam)
	st := ctx.URLParam(utils.StartTimeQueryParam)

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.ScenarioDetailsResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.ScenarioDetailsResponse

	if zkErr := validation.ValidateGetScenarioDetails(scenarioIds, st); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetScenarioDetailsHandler: ", zkErr)

	} else {
		resp, zkErr = t.service.GetScenarioDetailsService(scenarioIds, services, st)
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.ScenarioDetailsResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.ScenarioDetailsResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIssueDetailsHandler(ctx iris.Context) {
	issueHash := ctx.Params().Get(utils.IssueHash)

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IssueDetailsResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.IssueDetailsResponse

	if zkErr := validation.ValidateIssueDetailsHandler(issueHash); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetIssueDetailsHandler: ", zkErr)
	} else {
		resp, zkErr = t.service.GetIssueDetailsService(issueHash)
	}

	if zkErr == nil && resp.Issues.TotalCount == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IssueDetailsResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IssueDetailsResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetIncidentListHandler(ctx iris.Context) {
	issueHash := ctx.Params().Get(utils.IssueHash)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IncidentIdListResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.IncidentIdListResponse

	if zkErr := validation.ValidateIssueHashOffsetAndLimit(issueHash, offset, limit); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetIncidentListHandler api", zkErr)
	} else {
		l, _ := strconv.Atoi(limit)
		o, _ := strconv.Atoi(offset)
		resp, zkErr = t.service.GetIncidentListService(issueHash, o, l)
	}

	if zkErr == nil && resp.TotalRecords == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentIdListResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentIdListResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetPodDetailsHandler(ctx iris.Context) {

	// We can remove it from the url path. Trace_id is enough to identify everything. I'll add a to do as this would also require frontend changes.
	// TODO: The url path has issueHash, but we are not reading it here, remove from here and frontend.
	traceId := ctx.Params().Get(utils.IncidentId)
	spanId := ctx.URLParam(utils.SpanIdQueryParam)
	limit := ctx.URLParamDefault(utils.LimitQueryParam, "50")
	offset := ctx.URLParamDefault(utils.OffsetQueryParam, "0")

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.IncidentDetailsResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.IncidentDetailsResponse

	if zkErr := validation.ValidateGetIncidentDetailsApi(traceId, offset, limit); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetPodDetailsHandler api", zkErr)
	} else {
		l, _ := strconv.Atoi(limit)
		o, _ := strconv.Atoi(offset)
		resp, zkErr = t.service.GetIncidentDetailsService(traceId, spanId, o, l)
	}

	if zkErr == nil && len(resp.Spans) == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	// DONE
	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailsResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.IncidentDetailsResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetSpanRawDataHandler(ctx iris.Context) {
	// We can remove it from the url path. Trace_id is enough to identify everything. I'll add a to do as this would also require frontend changes.
	// TODO: The url path has issueHash, but we are not reading it here, remove from here and frontend.
	traceId := ctx.Params().Get(utils.IncidentId)
	spanId := ctx.Params().Get(utils.SpanId)

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.SpanRawDataResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.SpanRawDataResponse

	if zkErr := validation.ValidateGetSpanRawDataApi(traceId, spanId); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetSpanRawDataHandler api", zkErr)
	} else {
		resp, zkErr = t.service.GetSpanRawDataService(traceId, spanId)
	}

	if zkErr == nil && len(resp.Spans) == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.SpanRawDataResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.SpanRawDataResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func (t tracePersistenceHandler) GetErrorDataHandler(ctx iris.Context) {
	var errorReq request.GetErrorRequest
	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error reading request body")
		return
	}

	err = json.Unmarshal(body, &errorReq)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Error decoding JSON")
		return
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[traceResponse.ErrorDataResponse]
	var zkErr *zkErrorsUtils.ZkError
	var resp traceResponse.ErrorDataResponse

	if zkErr = validation.ValidateGetErrors(errorReq.ErrorIds); zkErr != nil {
		zkLogger.Error(LogTag, "Error while validating GetErrorHandler api", zkErr)
	} else {
		sanitizedErrorIds := make([]string, 0)
		for _, errorId := range errorReq.ErrorIds {
			if !zkCommon.IsEmpty(errorId) {
				sanitizedErrorIds = append(sanitizedErrorIds, errorId)
			}
		}
		if len(sanitizedErrorIds) == 0 {
			zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestErrorIdListIdEmpty, nil))
			zkLogger.Error(LogTag, "Error while validating GetErrorHandler api", zkErr)
			return
		} else {
			resp, zkErr = t.service.GetErrorDataService(errorReq.ErrorIds)
		}
	}

	if zkErr == nil && len(resp.Errors) == 0 {
		zkErr = zkCommon.ToPtr(zkErrorsUtils.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorNotFound, nil))
	}

	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.ErrorDataResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[traceResponse.ErrorDataResponse](200, resp, nil, zkErr)
	}

	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
