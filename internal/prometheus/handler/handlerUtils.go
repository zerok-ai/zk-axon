package handler

import (
	"axon/internal/prometheus/model/request"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
	"time"
)

func generatePromRequestMetadata(ctx iris.Context) request.PromRequestMeta {
	// Calculate the start and end times for the time range
	endTimeQP := ctx.URLParamDefault(utils.TimeQueryParam, time.Now().Format(time.Second.String()))
	durationQP := ctx.URLParamDefault(utils.DurationQueryParam, "10m")
	intervalQP := ctx.URLParamDefault(utils.RateIntervalQueryParam, "1m")

	endTime, err := time.Parse(time.Second.String(), endTimeQP)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing time: ", err)
		ctx.StatusCode(500)
		return request.PromRequestMeta{}
	}

	duration, err := time.ParseDuration(durationQP)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing duration: ", err)
		ctx.StatusCode(500)
		return request.PromRequestMeta{}
	}
	if duration > 0 {
		zkLogger.Error(LogTag, "Duration should be negative")
		ctx.StatusCode(500)
		return request.PromRequestMeta{}
	}

	interval, err := time.ParseDuration(intervalQP)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing interval: ", err)
		ctx.StatusCode(500)
		return request.PromRequestMeta{}
	}
	startTime := endTime.Add(duration)

	promReqMeta := request.PromRequestMeta{
		Namespace:    ctx.Params().Get(utils.Namespace),
		Pod:          ctx.Params().Get(utils.PodId),
		RateInterval: interval,
		StartTime:    startTime,
		EndTime:      endTime,
		Timestamp:    endTime.Unix(),
	}
	return promReqMeta
}

func generateGenericPromRequest(ctx iris.Context, req request.GenericHTTPRequest) *request.GenericPromRequest {
	promQuery := req.Query
	endTime := req.Time
	if endTime == 0 {
		endTime = time.Now().Unix()
	}

	durationStr := req.Duration
	if durationStr == "" {
		durationStr = "0m"
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing duration: ", err)
		ctx.StatusCode(500)
		return nil
	}

	stepStr := req.Step
	if stepStr == "" {
		stepStr = "1m"
	}
	step, err := time.ParseDuration(stepStr)
	if err != nil {
		zkLogger.Error(LogTag, "Error while parsing step: ", err)
		ctx.StatusCode(500)
		return nil
	}
	if step < 0 {
		zkLogger.Error(LogTag, "Step should be positive")
		ctx.StatusCode(500)
		return nil
	}

	startTime := endTime + int64(duration.Seconds())
	promIntegrationId := ctx.Params().Get(utils.IntegrationId)
	if promIntegrationId == "" {
		zkLogger.Error(LogTag, "Integration id is missing")
		ctx.StatusCode(500)
		return nil
	}
	queryReq := request.GenericPromRequest{
		PromIntegrationId: promIntegrationId,
		Query:             string(promQuery),
		StartTime:         int64(startTime),
		EndTime:           int64(endTime),
		Duration:          int64(duration),
		Step:              int64(step),
	}
	return &queryReq
}
func sendResponse[T any](ctx iris.Context, resp T, zkHttpResponse zkHttp.ZkHttpResponse[T], zkErr *zkErrors.ZkError, debug bool) {
	if debug {
		zkHttpResponse = zkHttp.ToZkResponse[T](200, resp, resp, zkErr)
	} else {
		zkHttpResponse = zkHttp.ToZkResponse[T](200, resp, nil, zkErr)
	}
	ctx.StatusCode(zkHttpResponse.Status)
	ctx.JSON(zkHttpResponse)
}

func arrayToPromList(arr []string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(strings.Join(arr, "|"), "|"), "|"))
}
