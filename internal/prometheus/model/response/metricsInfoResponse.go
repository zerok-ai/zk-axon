package response

import "time"

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

type LabelNameResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type MetricAttributesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    time.Time         `json:"activeAt"`
	Value       string            `json:"value"`
}

type Data struct {
	Alerts []Alert `json:"alerts"`
}

type AlertResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type IntegrationAlertsListResponse struct {
	Alerts []string `json:"alerts"`
	ErrorField
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
