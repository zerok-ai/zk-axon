package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"time"

	"github.com/lib/pq"
	zkcommon "github.com/zerok-ai/zk-utils-go/common"
)

type IncidentDetailsResponse struct {
	Spans SpansMetadataDetailsMap `json:"spans"`
}

type SpansMetadataDetailsMap map[string]SpanDetails

type SpanDetails struct {
	Error               bool           `json:"error"`
	TraceID             string         `json:"trace_id"`
	ParentSpanID        string         `json:"parent_span_id"`
	SpanID              string         `json:"span_id"`
	SpanName            string         `json:"span_name"`
	IsRoot              bool           `json:"is_root"`
	Kind                string         `json:"kind"`
	StartTime           time.Time      `json:"start_time"`
	Latency             uint64         `json:"latency"`
	Source              string         `json:"source"`
	Destination         string         `json:"destination"`
	WorkloadIDList      pq.StringArray `json:"workload_id_list"`
	Protocol            string         `json:"protocol"`
	IssueHashList       pq.StringArray `json:"issue_hash_list"`
	RequestPayloadSize  uint64         `json:"request_payload_size"`
	ResponsePayloadSize uint64         `json:"response_payload_size"`
	Method              string         `json:"method"`
	Route               string         `json:"route"`
	Scheme              string         `json:"scheme"`
	Path                string         `json:"path"`
	Query               string         `json:"query"`
	Status              *float64       `json:"status"`
	Metadata            *string        `json:"metadata,omitempty"`
	Username            string         `json:"username"`
	SourceIP            string         `json:"source_ip"`
	DestinationIP       string         `json:"destination_ip"`
	ServiceName         string         `json:"service_name"`
	Errors              string         `json:"errors"`

	SpanAttributes     zkcommon.GenericMap `json:"span_attributes,omitempty"`
	ResourceAttributes zkcommon.GenericMap `json:"resource_attributes,omitempty"`
	ScopeAttributes    zkcommon.GenericMap `json:"scope_attributes,omitempty"`
}

func ConvertSpanToIncidentDetailsResponse(t []dto.SpanTableDto) *IncidentDetailsResponse {
	respMap := make(map[string]SpanDetails)
	for _, v := range t {
		s := SpanDetails{
			Error:               v.WorkloadIDList != nil || len(v.WorkloadIDList) != 0,
			TraceID:             v.TraceID,
			ParentSpanID:        v.ParentSpanID,
			SpanID:              v.SpanID,
			SpanName:            v.SpanName,
			IsRoot:              v.IsRoot,
			Kind:                v.Kind,
			StartTime:           v.StartTime,
			Latency:             v.Latency,
			Source:              v.Source,
			Destination:         v.Destination,
			WorkloadIDList:      v.WorkloadIDList,
			Protocol:            v.Protocol,
			IssueHashList:       v.IssueHashList,
			RequestPayloadSize:  v.RequestPayloadSize,
			ResponsePayloadSize: v.ResponsePayloadSize,
			Method:              v.Method,
			Route:               v.Route,
			Scheme:              v.Scheme,
			Path:                v.Path,
			Query:               v.Query,
			Status:              v.Status,
			Username:            v.Username,
			SourceIP:            v.SourceIP,
			DestinationIP:       v.DestinationIP,
			ServiceName:         v.ServiceName,
			Errors:              v.Errors,
			SpanAttributes:      v.SpanAttributes,
			ResourceAttributes:  v.ResourceAttributes,
			ScopeAttributes:     v.ScopeAttributes,
		}

		respMap[v.SpanID] = s
	}

	resp := IncidentDetailsResponse{Spans: respMap}

	return &resp
}
