package response

type ContainerInfoResponse struct {
	ContainerInfo VectorList `json:"container_info"`
}

type ContainerInfo struct {
	Container   string `json:"container"`
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
	ImageID     string `json:"image_id"`
	ImageSpec   string `json:"image_spec"`
}

type AttributesMap map[string]string
type VectorList []AttributesMap
