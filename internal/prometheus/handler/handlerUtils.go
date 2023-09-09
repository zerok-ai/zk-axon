package handler

import (
	"axon/internal/prometheus/model/request"
	"axon/utils"
	"github.com/kataras/iris/v12"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
	"time"
)

func generatePromRequest(ctx iris.Context) request.PromRequestMeta {
	// Calculate the start and end times for the time range
	// TODO: Replace with your specified time
	endTime := time.Now()         // Replace with your specified time
	timeRange := 10 * time.Minute // Time range around the specified time
	startTime := endTime.Add(-timeRange)

	// TODO: Replace with request params
	logger.Debug(LogTag, "store len:", ctx.Params().Store.Len())

	promReqMeta := request.PromRequestMeta{
		Namespace:    ctx.Params().Get(utils.Namespace),
		Pod:          ctx.Params().Get(utils.PodId),
		RateInterval: ctx.URLParamDefault(utils.RateInterval, "10m"),
		StartTime:    startTime,
		EndTime:      endTime,
		Timestamp:    endTime.Unix(),
	}
	return promReqMeta
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
