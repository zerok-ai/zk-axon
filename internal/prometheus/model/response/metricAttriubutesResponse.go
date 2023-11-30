package response

type MetricAttributesListResponse struct {
	Attributes map[string]int `json:"attributes"`
}

type MetricAttributesPrometheusResponse struct {
	Status string          `json:"status"`
	Data   []AttributesMap `json:"data"`
}
