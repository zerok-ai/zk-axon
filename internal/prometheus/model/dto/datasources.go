package dto

type Datasource struct {
	Id             string         `json:"id"`
	Type           DataSourceType `json:"type"`
	Url            string         `json:"url"`
	Authentication Authentication `json:"authentication"`
	Level          string         `json:"level"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	Deleted        bool           `json:"deleted"`
	Disabled       bool           `json:"disabled"`
}

type Authentication struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// DataSourceType could contain values of type "PROMETHEUS", "DATADOG"
type DataSourceType string

const (
	Prometheus DataSourceType = "PROMETHEUS"
	Datadog    DataSourceType = "DATADOG"
)

// Level could be of type CLUSTER or ORGANIZATION
type Level string

const (
	Cluster      Level = "CLUSTER"
	Organization Level = "ORGANIZATION"
)
