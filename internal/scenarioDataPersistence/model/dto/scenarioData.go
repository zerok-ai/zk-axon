package dto

import (
	"github.com/lib/pq"
	"time"
)

type IncidentTableDto struct {
	TraceId                string    `json:"trace_id"`
	IssueHash              string    `json:"issue_hash"`
	IncidentCollectionTime time.Time `json:"incident_collection_time"`
}

type SpanTableDto struct {
	TraceId        string         `json:"trace_id"`
	SpanId         string         `json:"span_id"`
	ParentSpanId   string         `json:"parent_span_id"`
	Source         string         `json:"source"`
	Destination    string         `json:"destination"`
	WorkloadIdList pq.StringArray `json:"workload_id_list"`
	Status         string         `json:"status"`
	Metadata       string         `json:"metadata"`
	LatencyNs      float32        `json:"latency_ns"`
	Protocol       string         `json:"protocol"`
	Time           *time.Time     `json:"time"`
}

type IssueDetailsDto struct {
	IssueHash       string         `json:"issue_hash"`
	IssueTitle      string         `json:"issue_title"`
	ScenarioId      string         `json:"scenario_id"`
	ScenarioVersion string         `json:"scenario_version"`
	Sources         pq.StringArray `json:"sources"`
	Destinations    pq.StringArray `json:"destinations"`
	TotalCount      int            `json:"total_count"`
	Velocity        float32        `json:"velocity"`
	FirstSeen       time.Time      `json:"first_seen"`
	LastSeen        time.Time      `json:"last_seen"`
	Incidents       pq.StringArray `json:"incidents"`
}

type SpanRawDataDetailsDto struct {
	TraceId         string `json:"trace_id"`
	SpanId          string `json:"span_id"`
	Protocol        string `json:"protocol"`
	RequestPayload  []byte `json:"request_payload"`
	ResponsePayload []byte `json:"response_payload"`
}
