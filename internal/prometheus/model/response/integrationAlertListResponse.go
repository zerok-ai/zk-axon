package response

type IntegrationAlertsListResponse struct {
	AlertList struct {
		Groups []struct {
			Name     string `json:"name"`
			File     string `json:"file"`
			Rules    []Rule `json:"rules"`
			Interval int    `json:"interval"`
			Limit    int    `json:"limit"`
		} `json:"groups"`
	} `json:"alert_list"`
}

type AlertListPrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		Groups []struct {
			Name     string `json:"name"`
			File     string `json:"file"`
			Rules    []Rule `json:"rules"`
			Interval int    `json:"interval"`
			Limit    int    `json:"limit"`
		} `json:"groups"`
	} `json:"data"`
}

type Rule struct {
	State          string            `json:"state"`
	Name           string            `json:"name"`
	Query          string            `json:"query"`
	Duration       int               `json:"duration"`
	KeepFiringFor  int               `json:"keepFiringFor"`
	Labels         map[string]string `json:"labels"`
	Annotations    map[string]string `json:"annotations"`
	Alerts         []Alert           `json:"alerts"`
	Health         string            `json:"health"`
	EvaluationTime float64           `json:"evaluationTime"`
	LastEvaluation string            `json:"lastEvaluation"`
	Type           string            `json:"type"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    string            `json:"activeAt"`
	Value       string            `json:"value"`
}

func ConvertAlertListPrometheusResponseToIntegrationAlertsListResponse(alertPrometheusResponse AlertListPrometheusResponse) IntegrationAlertsListResponse {
	var integrationAlertsListResponse IntegrationAlertsListResponse
	integrationAlertsListResponse.AlertList = alertPrometheusResponse.Data
	return integrationAlertsListResponse
}
