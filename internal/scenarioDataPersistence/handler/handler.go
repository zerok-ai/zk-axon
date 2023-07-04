package handler

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/service"
	"axon/internal/scenarioDataPersistence/validation"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	"strconv"
)

type TracePersistenceHandler interface {
	GetIncidents(ctx iris.Context)
	GetTraces(ctx iris.Context)
	GetSpan(ctx iris.Context)
	GetSpanRawData(ctx iris.Context)
	GetMetadataMapData(ctx iris.Context)
}

type tracePersistenceHandler struct {
	service service.TracePersistenceService
}

func NewTracePersistenceHandler(persistenceService service.TracePersistenceService) TracePersistenceHandler {
	return &tracePersistenceHandler{
		service: persistenceService,
	}
}

// GetIncidents godoc
// @Summary Get Incidents
// @Description Get Incidents grouped by scenario title, destination and scenario type with pagination for given source
// @Tags trace persistence
// @Accept json
// @Produce json
// @Param scenarioType query string true "Scenario Type"
// @Param source query string true "Source"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} traceResponse.IncidentResponse
// @Router /trace-persistence/incidents [get]
func (t tracePersistenceHandler) GetIncidents(ctx iris.Context) {
	scenarioType := ctx.URLParam(utils.ScenarioType)
	source := ctx.URLParam(utils.Source)

	limit := ctx.URLParamDefault("limit", "50")
	offset := ctx.URLParamDefault("offset", "0")
	if err := validation.ValidateGetIncidentsDataApi(scenarioType, source, offset, limit); err != nil {
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetIncidentData(scenarioType, source, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.IncidentResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

// GetTraces godoc
// @Summary Get Traces
// @Description Get TraceId List for given scenario id with pagination
// @Tags trace persistence
// @Accept json
// @Produce json
// @Param scenarioId query string true "Scenario Id"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} traceResponse.TraceResponse
// @Router /trace-persistence/traces [get]
func (t tracePersistenceHandler) GetTraces(ctx iris.Context) {
	scenarioId := ctx.URLParam(utils.ScenarioId)
	limit := ctx.URLParamDefault(utils.Limit, "50")
	offset := ctx.URLParamDefault(utils.Offset, "0")
	if err := validation.ValidateGetTracesApi(scenarioId, offset, limit); err != nil {
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetTraces(scenarioId, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.TraceResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

// GetSpan godoc
// @Summary Get Span
// @Description Get Span details for given trace id, span id with pagination. If span id is not provided, it will return all spans
// @Tags trace persistence
// @Accept json
// @Produce json
// @Param traceId query string true "Trace Id"
// @Param spanId query string true "Span Id"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} traceResponse.SpanResponse
// @Router /trace-persistence/span [get]
func (t tracePersistenceHandler) GetSpan(ctx iris.Context) {
	traceId := ctx.URLParam(utils.TraceId)
	spanId := ctx.URLParam(utils.SpanId)
	limit := ctx.URLParamDefault(utils.Limit, "50")
	offset := ctx.URLParamDefault(utils.Offset, "0")
	if err := validation.ValidateGetTracesMetadataApi(traceId, offset, limit); err != nil {
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetTracesMetadata(traceId, spanId, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.SpanResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

// GetSpanRawData godoc
// @Summary Get raw data for given trace id, span id with pagination. Here raw data means request and response payload
// @Description Get raw data for given trace id, span id with pagination
// @Tags trace persistence
// @Accept json
// @Produce json
// @Param traceId query string true "Trace Id"
// @Param spanId query string true "Span Id"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} traceResponse.SpanRawDataResponse
// @Router /trace-persistence/span-raw-data [get]
func (t tracePersistenceHandler) GetSpanRawData(ctx iris.Context) {
	traceId := ctx.URLParam(utils.TraceId)
	spanId := ctx.URLParam(utils.SpanId)
	limit := ctx.URLParamDefault(utils.Limit, "50")
	offset := ctx.URLParamDefault(utils.Offset, "0")
	if err := validation.ValidateGetTracesRawDataApi(traceId, spanId, offset, limit); err != nil {
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetTracesRawData(traceId, spanId, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.SpanRawDataResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

// GetMetadataMapData godoc
// @Summary Get metadata map data
// @Description Get total trace count and list of protocols between all the sources and destinations which encountered an error, in given duration with pagination
// @Tags trace persistence
// @Accept json
// @Produce json
// @Param duration query string true "Duration"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {object} traceResponse.MetadataMapResponse
// @Router /trace-persistence/metadata-map [get]
func (t tracePersistenceHandler) GetMetadataMapData(ctx iris.Context) {
	d := ctx.URLParam(utils.Duration)
	limit := ctx.URLParamDefault(utils.Limit, "50")
	offset := ctx.URLParamDefault(utils.Offset, "0")
	if err := validation.ValidateGetMetadataMapApi(d, offset, limit); err != nil {
		z := &zkHttp.ZkHttpResponseBuilder[any]{}
		zkHttpResponse := z.WithZkErrorType(err.Error).Build()
		ctx.StatusCode(zkHttpResponse.Status)
		ctx.JSON(zkHttpResponse)
		return
	}

	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)

	resp, err := t.service.GetMetadataMap(d, o, l)

	zkHttpResponse := zkHttp.ToZkResponse[traceResponse.MetadataMapResponse](200, resp, resp, err)
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
