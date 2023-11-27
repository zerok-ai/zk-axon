package dto

import (
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type Integration struct {
	Id             string          `json:"id"`
	Alias          string          `json:"alias"`
	Type           IntegrationType `json:"type"`
	Url            string          `json:"url"`
	Authentication Authentication  `json:"authentication"`
	Level          string          `json:"level"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
	Deleted        bool            `json:"deleted"`
	Disabled       bool            `json:"disabled"`
	MetricServer   bool            `json:"metric_server"`
}

func (d Integration) Equals(other interfaces.ZKComparable) bool {
	return d.Id == other.(Integration).Id
}

type Authentication struct {
	Password *string `json:"password"`
	Username *string `json:"username"`
}

type UnsavedIntegrationRequestBody struct {
	Url string `json:"url"`
	Authentication
}

// IntegrationType could contain values of type "PROMETHEUS", "DATADOG"
type IntegrationType string

const (
	PrometheusIntegrationType IntegrationType = "PROMETHEUS"
	DatadogIntegrationType    IntegrationType = "DATADOG"
)

// Level could be of type CLUSTER or ORGANIZATION
type Level string

const (
	ClusterIntegrationLevel      Level = "CLUSTER"
	OrganizationIntegrationLevel Level = "ORGANIZATION"
)
