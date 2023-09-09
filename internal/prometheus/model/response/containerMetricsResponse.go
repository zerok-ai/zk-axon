package response

type ContainerMetricsResponse struct {
	CPUUsage Usage `json:"cpu_usage"`
	MemUsage Usage `json:"mem_usage"`
}

type PlotValues struct {
	TimeStamp []int64   `json:"time_stamp"`
	Values    []float64 `json:"values"`
}

type ContainerName string
type Usage map[ContainerName]PlotValues
