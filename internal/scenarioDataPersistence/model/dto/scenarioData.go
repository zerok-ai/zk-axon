package dto

import (
	"github.com/lib/pq"
	"time"
)

type IncidentTableDto struct {
	TotalRows              int       `json:"total_rows"`
	TraceId                string    `json:"trace_id"`
	IssueHash              string    `json:"issue_hash"`
	IncidentCollectionTime time.Time `json:"incident_collection_time"`
	EntryService           string    `json:"entry_service"`
	EndPoint               string    `json:"end_point"`
	Protocol               string    `json:"protocol"`
	RootSpanTime           time.Time `json:"root_span_time"`
	LatencyNs              *float32  `json:"latency_ns"`
}

type SpanTableDto struct {
	TraceID             string         `json:"trace_id"`
	ParentSpanID        string         `json:"parent_span_id"`
	SpanID              string         `json:"span_id"`
	IsRoot              bool           `json:"is_root"`
	Kind                string         `json:"kind"`
	StartTime           time.Time      `json:"start_time"`
	Latency             float32        `json:"latency"`
	Source              string         `json:"source"`
	Destination         string         `json:"destination"`
	WorkloadIDList      pq.StringArray `json:"workload_id_list"`
	Protocol            string         `json:"protocol"`
	IssueHashList       pq.StringArray `json:"issue_hash_list"`
	RequestPayloadSize  int64          `json:"request_payload_size"`
	ResponsePayloadSize int64          `json:"response_payload_size"`
	Method              string         `json:"method"`
	Route               string         `json:"route"`
	Scheme              string         `json:"scheme"`
	Path                string         `json:"path"`
	Query               string         `json:"query"`
	Status              int            `json:"status"`
	Metadata            *string        `json:"metadata"`
	Username            string         `json:"username"`
	SourceIP            string         `json:"source_ip"`
	DestinationIP       string         `json:"destination_ip"`
	ServiceName         string         `json:"service_name"`
	ErrorType           string         `json:"error_type"`
	ErrorTableId        string         `json:"error_table_id"`
}

type IssueDetailsDto struct {
	TotalRows       int            `json:"total_rows"`
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

type ScenarioDetailsDto struct {
	ScenarioId      string         `json:"scenario_id"`
	ScenarioVersion string         `json:"scenario_version"`
	Sources         pq.StringArray `json:"sources"`
	Destinations    pq.StringArray `json:"destinations"`
	TotalCount      int            `json:"total_count"`
	Velocity        float32        `json:"velocity"`
	FirstSeen       time.Time      `json:"first_seen"`
	LastSeen        time.Time      `json:"last_seen"`
}

type SpanRawDataDetailsDto struct {
	TraceID     string `json:"trace_id"`
	SpanID      string `json:"span_id"`
	Protocol    string `json:"protocol"`
	ReqHeaders  string `json:"req_headers"`
	RespHeaders string `json:"resp_headers"`
	IsTruncated bool   `json:"is_truncated"`
	ReqBody     []byte `json:"req_body"`
	RespBody    []byte `json:"resp_body"`
}

type ErrorDataTableDto struct {
	Id   string `json:"id"`
	Data []byte `json:"data"`
}
