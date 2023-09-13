package dto

import (
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type Datasource struct {
	Id             string         `json:"id"`
	Alias          string         `json:"alias"`
	Type           DataSourceType `json:"type"`
	Url            string         `json:"url"`
	Authentication Authentication `json:"authentication"`
	Level          string         `json:"level"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	Deleted        bool           `json:"deleted"`
	Disabled       bool           `json:"disabled"`
	MetricServer   bool           `json:"metric_server"`
}

func (d Datasource) Equals(other interfaces.ZKComparable) bool {
	return d.Id == other.(Datasource).Id
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
