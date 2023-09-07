package response

import (
	"github.com/prometheus/common/model"
)

type PodDetailResponse struct {
	PodName      string   `json:"pod_name"`
	Metadata     Metadata `json:"metadata"`
	ZkInferences string   `json:"zk_inferences"`
	CPUUsage     Usage    `json:"cpu_usage"`
	MemUsage     Usage    `json:"mem_usage"`
}

type Usage struct {
	Title   string  `json:"title"`
	Success bool    `json:"success"`
	Frames  []Frame `json:"frames"`
}

type Frame struct {
	Schema Schema `json:"schema"`
	Data   Data   `json:"data"`
}

type Data struct {
	TimeStamp []int64   `json:"time_stamp"`
	Values    []float64 `json:"values"`
}

type Schema struct {
	Name string `json:"name"`
}

type Metadata interface {
}

type Metadata_old struct {
	AppKubernetesIoComponent string `json:"app_kubernetes_io_component"`
	AppKubernetesIoInstance  string `json:"app_kubernetes_io_instance"`
	AppKubernetesIoManagedBy string `json:"app_kubernetes_io_managed_by"`
	AppKubernetesIoName      string `json:"app_kubernetes_io_name"`
	AppKubernetesIoPartOf    string `json:"app_kubernetes_io_part_of"`
	AppKubernetesIoVersion   string `json:"app_kubernetes_io_version"`
	Container                string `json:"container"`
	ContainerID              string `json:"container_id"`
	CreatedByKind            string `json:"created_by_kind"`
	CreatedByName            string `json:"created_by_name"`
	HelmShChart              string `json:"helm_sh_chart"`
	HostIP                   string `json:"host_ip"`
	HostNetwork              string `json:"host_network"`
	Image                    string `json:"image"`
	ImageID                  string `json:"image_id"`
	ImageSpec                string `json:"image_spec"`
	Instance                 string `json:"instance"`
	Job                      string `json:"job"`
	Namespace                string `json:"namespace"`
	Node                     string `json:"node"`
	Pod                      string `json:"pod"`
	PodIP                    string `json:"pod_ip"`
	Service                  string `json:"service"`
	Uid                      string `json:"uid"`
}

func ConvertMetricToPodUsage(title string, metric model.Matrix) Usage {
	var usage Usage
	var frames []Frame
	var data Data
	var schema Schema

	for _, series := range metric {
		var values []float64
		var timeStamps []int64

		for _, value := range series.Values {
			values = append(values, float64(value.Value))
			timeStamps = append(timeStamps, int64(value.Timestamp))
		}

		data.Values = values
		data.TimeStamp = timeStamps

		// Add container name
		var labelSet = make(map[string]string)
		for key, value := range series.Metric {
			println(string(key), string(value))
			labelSet[string(key)] = string(value)
		}
		println(labelSet["container"])
		schema.Name = labelSet["container"]
		// Skip adding empty container names
		if schema.Name == "" {
			continue
		}

		frame := Frame{
			Schema: schema,
			Data:   data,
		}

		frames = append(frames, frame)
	}

	usage.Title = title
	usage.Success = true
	usage.Frames = frames

	return usage
}
