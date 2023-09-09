package prometheus

import (
	"axon/internal/prometheus/handler"
	"axon/utils"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.PrometheusHandler) {

	promAPI := app.Party("/c/axon/prom")
	{
		promAPI.Get("/pods-info/trace/{"+utils.TraceId+"}", tph.GetPodsInfoHandler)
		promAPI.Get("/container-info/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", tph.GetContainerInfoHandler)
		promAPI.Get("/container-metrics/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", tph.GetContainerMetricsHandler)
	}
}
