package dto

import (
	"github.com/lib/pq"
	"time"
)

var LogTag = "trace_dto"

type IncidentDto struct {
	ScenarioId      string  `json:"scenario_id"`
	ScenarioVersion string  `json:"scenario_version"`
	Title           string  `json:"title"`
	ScenarioType    string  `json:"scenario_type"`
	Velocity        float32 `json:"velocity"`
	TotalCount      int     `json:"total_count"`
	Source          string  `json:"source"`
	Destination     string  `json:"destination"`
	FirstSeen       string  `json:"first_seen"`
	LastSeen        string  `json:"last_seen"`
}

type ScenarioTableDto struct {
	ScenarioId      string    `json:"scenario_id"`
	ScenarioVersion string    `json:"scenario_version"`
	TraceId         string    `json:"trace_id"`
	ScenarioTitle   string    `json:"scenario_title"`
	ScenarioType    string    `json:"scenario_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type SpanTableDto struct {
	TraceId        string         `json:"trace_id"`
	SpanId         string         `json:"span_id"`
	ParentSpanId   string         `json:"parent_span_id"`
	Source         string         `json:"source"`
	Destination    string         `json:"destination"`
	WorkloadIdList pq.StringArray `json:"workload_id_list"`
	Metadata       string         `json:"metadata"`
	LatencyMs      float32        `json:"latency_ms"`
	Protocol       string         `json:"protocol"`
}

type SpanRawDataTableDto struct {
	TraceId         string `json:"trace_id"`
	SpanId          string `json:"span_id"`
	RequestPayload  []byte `json:"request_payload"`
	ResponsePayload []byte `json:"response_payload"`
}

type MetadataMapDto struct {
	Source       string         `json:"source"`
	Destination  string         `json:"destination"`
	TraceCount   int            `json:"trace_count"`
	ProtocolList pq.StringArray `json:"protocol_list"`
}
