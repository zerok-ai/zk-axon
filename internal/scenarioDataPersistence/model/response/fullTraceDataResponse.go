package scenariodataresponse

//
//import (
//	"axon/internal/scenarioDataPersistence/model/dto"
//	"database/sql/driver"
//	"encoding/json"
//	"errors"
//	"time"
//)
//
//type ScenarioIncidentDetailsResponse struct {
//	ScenarioDetails ScenariosIncidentDetails `json:"scenario_details"`
//}
//
//func ConvertScenarioIncidentDetailsDtoListToScenarioIncidentDetailsResponse(d []dto.ScenarioIncidentDetailsDto) ScenarioIncidentDetailsResponse {
//	respmap := make(map[string]map[string][]Span, 0)
//
//	for _, v := range d {
//		if _, ok := respmap[v.ScenarioId]; !ok {
//			respmap[v.ScenarioId] = make(map[string][]Span, 0)
//		}
//		if _, ok := respmap[v.ScenarioId][v.TraceId]; !ok {
//			respmap[v.ScenarioId][v.TraceId] = make([]Span, 0)
//		}
//
//		s := Span{
//			SpanId:          v.SpanId,
//			TraceId:         v.TraceId,
//			ParentSpanId:    v.ParentSpanId,
//			Source:          v.Source,
//			Destination:     v.Destination,
//			WorkloadIdList:  v.WorkloadIdList,
//			Metadata:        v.Metadata,
//			LatencyMs:       v.LatencyMs,
//			Protocol:        v.Protocol,
//			RequestPayload:  HTTPRequestPayload{},
//			ResponsePayload: HTTPResponsePayload{},
//			IssueHashList:   v.IssueHashList,
//			Time:            v.Time,
//		}
//
//		respmap[v.ScenarioId][v.TraceId] = append(respmap[v.ScenarioId][v.TraceId], s)
//	}
//	return ScenarioIncidentDetailsResponse{}
//}
//
//type ScenariosIncidentDetails struct {
//	ScenarioId          string     `json:"scenario_id"`
//	IncidentDetailsList []Incident `json:"incident_details_list"`
//}
//
//type Incident struct {
//	TraceId                string    `json:"trace_id"`
//	Spans                  []Span    `json:"spans"`
//	IncidentCollectionTime time.Time `json:"incident_collection_time"`
//}
//
//type Span struct {
//	SpanId          string              `json:"span_id"`
//	TraceId         string              `json:"trace_id"`
//	ParentSpanId    string              `json:"parent_span_id"`
//	Source          string              `json:"source"`
//	Destination     string              `json:"destination"`
//	WorkloadIdList  []string            `json:"workload_id_list"`
//	Metadata        Metadata            `json:"metadata"`
//	LatencyMs       float32             `json:"latency_ms"`
//	Protocol        string              `json:"protocol"`
//	RequestPayload  HTTPRequestPayload  `json:"request_payload"`
//	ResponsePayload HTTPResponsePayload `json:"response_payload"`
//	IssueHashList   []string            `json:"issue_hash_list"`
//	Time            *time.Time          `json:"time"`
//}
//
//type Metadata map[string]interface{}
//
//type HTTPRequestPayload struct {
//	ReqPath    string `json:"req_path"`
//	ReqMethod  string `json:"req_method"`
//	ReqHeaders string `json:"req_headers"`
//	ReqBody    string `json:"req_body"`
//}
//
//type HTTPResponsePayload struct {
//	RespStatus  string `json:"resp_status"`
//	RespMessage string `json:"resp_message"`
//	RespHeaders string `json:"resp_headers"`
//	RespBody    string `json:"resp_body"`
//}
//
//// Value Make the Attrs struct implement the driver.Valuer interface. This method
//// simply returns the JSON-encoded representation of the struct.
//func (a Metadata) Value() (driver.Value, error) {
//	return json.Marshal(a)
//}
//
//// Scan Make the Attrs struct implement the sql.Scanner interface. This method
//// simply decodes a JSON-encoded value into the struct fields.
//func (a *Metadata) Scan(value interface{}) error {
//	b, ok := value.([]byte)
//	if !ok {
//		return errors.New("type assertion to []byte failed")
//	}
//
//	return json.Unmarshal(b, &a)
//}
