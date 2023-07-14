package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"github.com/lib/pq"
	"time"
)

type IncidentDetailsResponse struct {
	Spans SpansMetadataDetailsMap `json:"spans"`
}

type SpansMetadataDetailsMap map[string]SpanDetails

type SpanDetails struct {
	Source         string         `json:"source"`
	Destination    string         `json:"destination"`
	Error          bool           `json:"error"`
	Metadata       string         `json:"metadata,omitempty"`
	LatencyNs      float32        `json:"latency_ns"`
	Protocol       string         `json:"protocol"`
	Status         string         `json:"status"`
	ParentSpanId   string         `json:"parent_span_id"`
	WorkloadIdList pq.StringArray `json:"workload_id_list"`
	Time           *time.Time     `json:"time"`
}

func ConvertSpanToIncidentDetailsResponse(t []dto.SpanTableDto) (*IncidentDetailsResponse, *error) {
	respMap := make(map[string]SpanDetails, 0)
	for _, v := range t {

		s := SpanDetails{
			Source:         v.Source,
			Destination:    v.Destination,
			Error:          v.WorkloadIdList != nil || len(v.WorkloadIdList) != 0,
			Metadata:       v.Metadata,
			LatencyNs:      v.LatencyNs,
			Protocol:       v.Protocol,
			Status:         v.Status,
			ParentSpanId:   v.ParentSpanId,
			WorkloadIdList: v.WorkloadIdList,
			Time:           v.Time,
		}

		respMap[v.SpanId] = s
	}

	resp := IncidentDetailsResponse{Spans: respMap}

	return &resp, nil
}
