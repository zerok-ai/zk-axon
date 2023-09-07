package handler

import (
	"axon/internal/config"
	"axon/internal/prometheus/model/request"
	promResponse "axon/internal/prometheus/model/response"
	"axon/internal/prometheus/service"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	"time"
)

type PrometheusHandler interface {
	GetPodDetailsHandler(ctx iris.Context)
}

var LogTag = "prometheus_handler"

type prometheusHandler struct {
	service service.PrometheusService
	cfg     config.AppConfigs
}

func NewPrometheusHandler(persistenceService service.PrometheusService, cfg config.AppConfigs) PrometheusHandler {
	return &prometheusHandler{
		service: persistenceService,
		cfg:     cfg,
	}
}

func (t prometheusHandler) GetPodDetailsHandler(ctx iris.Context) {
	// Calculate the start and end times for the time range
	endTime := time.Now()         // Replace with your specified time
	timeRange := 10 * time.Minute // Time range around the specified time
	startTime := endTime.Add(-timeRange)

	podInfoReq := request.PodInfoRequest{
		Namespace:    ctx.Params().Get(utils.Namespace),
		Pod:          ctx.Params().Get(utils.PodId),
		RateInterval: ctx.URLParamDefault(utils.RateInterval, "10m"),
		StartTime:    startTime,
		EndTime:      endTime,
		Timestamp:    endTime.Unix(),
	}

	var zkHttpResponse zkHttp.ZkHttpResponse[promResponse.PodDetailResponse]

	resp, zkErr := t.service.GetPodDetailsService(podInfoReq)
	if t.cfg.Http.Debug {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.PodDetailResponse](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[promResponse.PodDetailResponse](200, resp, nil, zkErr)
	}
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}
