package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"time"
)

type IncidentDetailListResponse struct {
	IncidentDetList []IncidentDetail `json:"trace_det_list"`
	TotalRecords    int              `json:"total_records"`
}

type IncidentDetail struct {
	IncidentId             string    `json:"incident_id"`
	EntryService           string    `json:"entry_service"`
	EndPoint               string    `json:"entry_path"`
	RootSpanTime           time.Time `json:"root_span_time"`
	LatencyNs              *float32  `json:"latency_ns"`
	IncidentCollectionTime time.Time `json:"incident_collection_time"`
	Protocol               string    `json:"protocol"`
}

func ConvertIncidentTableDtoToIncidentDetailListResponse(t []dto.IncidentTableDto) *IncidentDetailListResponse {
	incidentDetailList := make([]IncidentDetail, 0)
	for _, v := range t {
		incidentDetailList = append(incidentDetailList, getIncidentDetail(v))
	}

	if len(incidentDetailList) > 0 {
		return &IncidentDetailListResponse{IncidentDetList: incidentDetailList, TotalRecords: t[0].TotalRows}
	}

	return &IncidentDetailListResponse{IncidentDetList: incidentDetailList, TotalRecords: 0}
}

func getIncidentDetail(t dto.IncidentTableDto) IncidentDetail {
	return IncidentDetail{
		IncidentId:             t.TraceId,
		EntryService:           t.EntryService,
		EndPoint:               t.EndPoint,
		RootSpanTime:           t.RootSpanTime,
		LatencyNs:              t.LatencyNs,
		IncidentCollectionTime: t.IncidentCollectionTime,
		Protocol:               t.Protocol,
	}
}
