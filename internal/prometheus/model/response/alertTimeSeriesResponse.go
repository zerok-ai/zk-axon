package response

type AlertTimeSeriesResponse struct {
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

type AlertRangePrometheusResponse struct {
	Status string         `json:"status"`
	Data   AlertRangeData `json:"data"`
}

type AlertRangeData struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric Metric  `json:"metric"`
	Values []Value `json:"values"`
}

type Metric struct {
	Name       string `json:"__name__"`
	AlertName  string `json:"alertname"`
	AlertState string `json:"alertstate"`
	Pod        string `json:"pod"`
	Severity   string `json:"severity"`
}

type Value []interface{}
