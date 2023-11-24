package response

type IsIntegrationMetricServerResponse struct {
	MetricServer *bool `json:"metric_server,omitempty"`
	ErrorField
}

type MetricAttributesListResponse struct {
	Attributes map[string]int `json:"attributes"`
	ErrorField
}

type IntegrationMetricsListResponse struct {
	Metrics []string `json:"metrics"`
	ErrorField
}

type ErrorField struct {
	StatusCode *int    `json:"status_code,omitempty"`
	Status     *string `json:"status,omitempty"`
	Error      *bool   `json:"error,omitempty"`
}

type IntegrationAlertsListResponse struct {
	Alerts []string `json:"alerts"`
	ErrorField
}

type LabelNameResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type MetricAttributesResponse struct {
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
	ConnectionStatus  string `json:"connection_status"`
	ConnectionMessage string `json:"connection_message,omitempty"`
	HasMetricServer   *bool  `json:"has_metric_server,omitempty"`
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

type MetricAttributes struct {
	Status string          `json:"status"`
	Data   []AttributesMap `json:"data"`
}
