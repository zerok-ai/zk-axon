package response

import "time"

type IntegrationAlertsListResponse struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    time.Time         `json:"activeAt"`
	Value       string            `json:"value"`
}

type AlertPrometheusResponse struct {
	Data AlertsData `json:"data"`
}

type AlertsData struct {
	Alerts []Alert `json:"alerts"`
}
