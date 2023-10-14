package request

type GenericPromRequest struct {
	PromIntegrationId string `json:"prom_integration_id"`
	Query             string `json:"query"`
	StartTime         int64  `json:"start_time,omitempty"`
	EndTime           int64  `json:"end_time,omitempty"`
	Duration          int64  `json:"duration,omitempty"`
	Step              int64  `json:"step,omitempty"`
}

type GenericHTTPRequest struct {
	Query    string `json:"query"`
	Time     int64  `json:"time,omitempty"`
	Duration string `json:"duration,omitempty"`
	Step     string `json:"step,omitempty"`
}
