package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"axon/utils"
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
	IssueHash       string    `json:"issue_hash"`
	IssueTitle      string    `json:"issue_title"`
	ScenarioId      string    `json:"scenario_id"`
	ScenarioVersion string    `json:"scenario_version"`
	Sources         []string  `json:"sources"`
	Destinations    []string  `json:"destinations"`
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
	issuesList := make([]IssueDetails, 0)

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
	hours := utils.HoursBetween(v.FirstSeen, v.LastSeen) + 1

	r.IssueHash = v.IssueHash
	r.IssueTitle = v.IssueTitle
	r.ScenarioId = v.ScenarioId
	r.ScenarioVersion = v.ScenarioVersion
	r.Sources = v.Sources
	r.Destinations = v.Destinations
	r.TotalCount = v.TotalCount
	r.Velocity = float32(v.TotalCount / hours)
	r.FirstSeen = v.FirstSeen
	r.LastSeen = v.LastSeen
	if len(v.Incidents) >= 5 {
		r.Incidents = v.Incidents[:5]
	} else {
		r.Incidents = v.Incidents
	}

	if len(v.Sources) >= 5 {
		r.Sources = v.Sources[:5]
	} else {
		r.Sources = v.Sources
	}

	if len(v.Destinations) >= 5 {
		r.Destinations = v.Destinations[:5]
	} else {
		r.Destinations = v.Destinations
	}

	return r
}
