package response

type IsIntegrationMetricServerResponse struct {
	MetricServer bool `json:"metric_server"`
}

type IntegrationMetricsListResponse struct {
	Metrics []string `json:"metrics"`
}

type IntegrationAlertsListResponse struct {
	Alerts []string `json:"alerts"`
}

type LabelNameResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type AlertsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []interface{} `json:"alerts"`
	} `json:"data"`
}

type TestConnectionResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
}

type QueryResult struct {
	Status    string      `json:"status"`
	Data      DataSection `json:"data"`
	ErrorType string      `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}

// DataSection represents the data section of the query result
type DataSection struct {
	ResultType string      `json:"resultType"`
	Result     interface{} `json:"result"`
}
