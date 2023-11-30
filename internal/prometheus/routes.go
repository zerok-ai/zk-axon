package prometheus

import (
	"axon/internal/prometheus/handler"
	"axon/utils"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, ph handler.PrometheusHandler) {
	promClusterAPIs := app.Party("/c/axon")
	{
		promClusterAPIs.Get("/prom/pods-info/trace/{"+utils.TraceId+"}", ph.GetPodsInfoHandler)
		promClusterAPIs.Get("/prom/container-info/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", ph.GetContainerInfoHandler)
		promClusterAPIs.Get("/prom/container-metrics/pod/{"+utils.Namespace+"}/{"+utils.PodId+"}", ph.GetContainerMetricsHandler)

		promClusterAPIs.Post("/prom/{"+utils.IntegrationId+"}/query", ph.GetGenericQueryHandler)

		promClusterAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/status", ph.TestIntegrationConnectionStatus)
		promClusterAPIs.Post("/prom/unsaved/status", ph.TestUnsavedIntegrationConnectionStatus)
		//promClusterAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/metrics", ph.GetMetrics)                                                             //
		//promClusterAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/metric/{"+utils.MetricAttributeNamePathParam+"}/attributes", ph.GetMetricAttributes) //
		promClusterAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/alerts", ph.GetAlerts)
		promClusterAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/alerts/range", ph.GetAlertsTimeSeries)
	}

	promClusterInternalAPIs := app.Party("/i/axon")
	{
		promClusterInternalAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/metrics", ph.GetMetrics)
		promClusterInternalAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/metric/{"+utils.MetricAttributeNamePathParam+"}/attributes", ph.GetMetricAttributes)
		promClusterInternalAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/alerts", ph.GetAlerts)
		promClusterInternalAPIs.Get("/prom/{"+utils.IntegrationIdxPathParam+"}/alerts/range", ph.GetAlertsTimeSeries)
		promClusterInternalAPIs.Post("/prom/{"+utils.IntegrationIdxPathParam+"}/alerts/webhook", ph.PrometheusAlertWebhook)

	}

}
