package prometheus

import (
	"axon/internal/prometheus/handler"
	"axon/utils"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.PrometheusHandler) {
	promClusterAPIs := app.Party("/c/axon")
	{
		promClusterAPIs.Get("/prom/pods-info/trace/{"+utils.TraceId+"}", tph.GetPodsInfoHandler)
		promClusterAPIs.Get("/prom/container-info/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", tph.GetContainerInfoHandler)
		promClusterAPIs.Get("/prom/container-metrics/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", tph.GetContainerMetricsHandler)

		promClusterAPIs.Post("/prom/{"+utils.DatasourceId+"}/query", tph.GetGenericQueryHandler)
	}
}
