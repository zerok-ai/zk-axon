package response

type TestConnectionResponse struct {
	ConnectionStatus  string `json:"connection_status"`
	ConnectionMessage string `json:"connection_message,omitempty"`
	HasMetricServer   *bool  `json:"has_metric_server,omitempty"`
}
