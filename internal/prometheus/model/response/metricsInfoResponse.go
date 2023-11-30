package response

type IsIntegrationMetricServerResponse struct {
	MetricServer *bool `json:"metric_server,omitempty"`
	ErrorField
}

type ErrorField struct {
	StatusCode *int    `json:"status_code,omitempty"`
	Status     *string `json:"status,omitempty"`
	Error      *bool   `json:"error,omitempty"`
}

type StringListPrometheusResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type QueryResultPrometheusResponse struct {
	Status    string      `json:"status"`
	Data      DataSection `json:"data"`
	ErrorType string      `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}

type DataSection struct {
	ResultType string      `json:"resultType"`
	Result     interface{} `json:"result"`
}
