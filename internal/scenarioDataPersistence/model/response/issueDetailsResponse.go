package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"time"
)

var LogTag = "scenario_response"

type IncidentListResponse struct {
	TraceIdList []string `json:"trace_id_list"`
}

func ConvertIncidentTableDtoToIncidentListResponse(t []dto.IncidentTableDto) *IncidentListResponse {
	traceIdList := make([]string, 0)
	for _, v := range t {
		traceIdList = append(traceIdList, v.TraceId)
	}

	return &IncidentListResponse{TraceIdList: traceIdList}
}

type IssueDetails struct {
	IssueId         string    `json:"issue_id"`
	IssueTitle      string    `json:"issue_title"`
	ScenarioId      string    `json:"scenario_id"`
	ScenarioVersion string    `json:"scenario_version"`
	Source          string    `json:"source"`
	Destination     string    `json:"destination"`
	TotalCount      int       `json:"total_count"`
	Velocity        float32   `json:"velocity"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
	Incidents       []string  `json:"incidents"`
}

type IssueListWithDetailsResponse struct {
	Issues []IssueDetails `json:"issues"`
}

func ConvertIssueListDetailsDtoToIssueListDetailsResponse(t []dto.IssueDetailsDto) *IssueListWithDetailsResponse {
	var resp IssueListWithDetailsResponse
	var issuesList []IssueDetails

	for _, v := range t {
		r := ConvertIssueDetailsDtoToIssueDetails(v)
		issuesList = append(issuesList, r)
	}

	resp.Issues = issuesList

	return &resp
}

type IssueWithDetailsResponse struct {
	Issue IssueDetails `json:"issue"`
}

func ConvertIssueDetailsDtoToIssueDetails(v dto.IssueDetailsDto) IssueDetails {
	var r IssueDetails

	r.IssueId = v.IssueId
	r.IssueTitle = v.IssueTitle
	r.ScenarioId = v.ScenarioId
	r.ScenarioVersion = v.ScenarioVersion
	r.Source = v.Source
	r.Destination = v.Destination
	r.TotalCount = v.TotalCount
	r.Velocity = v.Velocity
	r.FirstSeen = v.FirstSeen
	r.LastSeen = v.LastSeen
	if len(v.Incidents) >= 5 {
		r.Incidents = v.Incidents[:5]
	} else {
		r.Incidents = v.Incidents
	}

	return r
}

func ConvertIssueToIssueDetailsResponse(t dto.IssueDetailsDto) IssueWithDetailsResponse {
	return IssueWithDetailsResponse{Issue: ConvertIssueDetailsDtoToIssueDetails(t)}
}
