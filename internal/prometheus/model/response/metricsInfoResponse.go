package response

import "time"

type IsIntegrationMetricServerResponse struct {
	MetricServer *bool `json:"metric_server,omitempty"`
	ErrorField
}

type MetricAttributesListResponse struct {
	Attributes map[string]int `json:"attributes"`
}

type IntegrationMetricsListResponse struct {
	Metrics []string `json:"metrics"`
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

type AlertResponseData struct {
	Alerts []Alert `json:"alerts"`
}

type AlertResponse struct {
	Data AlertResponseData `json:"data"`
}

type IntegrationAlertsListResponse struct {
	Alerts []Alert `json:"alerts"`
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

type Metric struct {
	Name       string `json:"__name__"`
	AlertName  string `json:"alertname"`
	AlertState string `json:"alertstate"`
	Pod        string `json:"pod"`
	Severity   string `json:"severity"`
}

type Value []interface{}

type Result struct {
	Metric Metric  `json:"metric"`
	Values []Value `json:"values"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type AlertRangeApiResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type AlertRangeResponse struct {
	AlertsRangeData []AlertsRangeData `json:"alerts_range_data"`
}

type AlertsRangeData struct {
	AlertName  string       `json:"alert_name"`
	SeriesData []SeriesData `json:"series_data"`
}

type SeriesData struct {
	State    string     `json:"state"`
	Duration []Duration `json:"duration"`
}

type Duration struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type AlterTrigger struct {
	State  string  `json:"state"`
	Period [][]int `json:"period"`
}
