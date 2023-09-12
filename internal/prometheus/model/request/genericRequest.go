package request

type GenericRequest struct {
	PromDatasourceId string `json:"prom_datasource_id"`
	Query            string `json:"query"`
	StartTime        int64  `json:"start_time,omitempty"`
	EndTime          int64  `json:"end_time,omitempty"`
	Duration         int64  `json:"duration,omitempty"`
}

type GenericHTTPRequest struct {
	PromDatasourceId string `json:"prom_datasource_id"`
	Query            string `json:"query"`
	Time             int64  `json:"time,omitempty"`
	Duration         string `json:"duration,omitempty"`
}
