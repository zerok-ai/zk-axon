package prometheus

import (
	"axon/internal/prometheus/handler"
	"axon/utils"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.PrometheusHandler) {

	promAPI := app.Party("/c/prom")
	{
		promAPI.Get("/pod/{"+utils.PodId+"}/ns/{"+utils.Namespace+"}", tph.GetPodDetailsHandler)
	}
}
