package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"axon/utils"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"time"
)

var LogTag = "scenario_response"

type IncidentIdListResponse struct {
	TraceIdList  []string `json:"trace_id_list"`
	TotalRecords int      `json:"total_records"`
}

func ConvertIncidentTableDtoToIncidentListResponse(t []dto.IncidentTableDto) *IncidentIdListResponse {
	traceIdList := make([]string, 0)
	for _, v := range t {
		traceIdList = append(traceIdList, v.TraceId)
	}

	if len(traceIdList) > 0 {
		return &IncidentIdListResponse{TraceIdList: traceIdList, TotalRecords: t[0].TotalRows}
	}

	return &IncidentIdListResponse{TraceIdList: traceIdList, TotalRecords: 0}
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

type ScenarioDetails struct {
	ScenarioId      string    `json:"scenario_id"`
	ScenarioVersion string    `json:"scenario_version"`
	Sources         []string  `json:"sources"`
	Destinations    []string  `json:"destinations"`
	TotalCount      int       `json:"total_count"`
	Velocity        float32   `json:"velocity"`
	FirstSeen       time.Time `json:"first_seen"`
	LastSeen        time.Time `json:"last_seen"`
}

type IssueListWithDetailsResponse struct {
	Issues       []IssueDetails `json:"issues"`
	TotalRecords int            `json:"total_records"`
}

type ScenarioDetailsResponse struct {
	Scenarios []ScenarioDetails `json:"scenarios"`
}

type IssueDetailsResponse struct {
	Issues IssueDetails `json:"issue"`
}

func ConvertIssueListDetailsDtoToIssueListDetailsResponse(t []dto.IssueDetailsDto) *IssueListWithDetailsResponse {
	var resp IssueListWithDetailsResponse
	issuesList := make([]IssueDetails, 0)

	for _, v := range t {
		r := ConvertIssueDetailsDtoToIssueDetails(v)
		issuesList = append(issuesList, r)
	}

	resp.Issues = issuesList
	if len(issuesList) > 0 {
		resp.TotalRecords = t[0].TotalRows
	}

	return &resp
}

func ConvertScenarioDetailsDtoToScenarioDetailsResponse(t []dto.ScenarioDetailsDto) *ScenarioDetailsResponse {
	var resp ScenarioDetailsResponse
	scenarioDetails := make([]ScenarioDetails, 0)

	for _, v := range t {
		r := ConvertScenarioDetailsDtoToScenarioDetailsDetailsResponse(v)
		scenarioDetails = append(scenarioDetails, r)
	}

	resp.Scenarios = scenarioDetails
	return &resp
}

func ConvertIssueDetailsDtoToIssueListDetailsResponse(t []dto.IssueDetailsDto) IssueDetailsResponse {
	var resp IssueDetailsResponse

	if t != nil && len(t) > 0 {
		resp.Issues = ConvertIssueDetailsDtoToIssueDetails(t[0])
	} else {
		return resp
	}

	if len(t) > 1 {
		zkLogger.Info(LogTag, "IssueDetailsDto has more than one record")
	}

	return resp
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

func ConvertScenarioDetailsDtoToScenarioDetailsDetailsResponse(v dto.ScenarioDetailsDto) ScenarioDetails {
	var r ScenarioDetails
	hours := utils.HoursBetween(v.FirstSeen, v.LastSeen) + 1

	r.ScenarioId = v.ScenarioId
	r.ScenarioVersion = v.ScenarioVersion
	r.Sources = v.Sources
	r.Destinations = v.Destinations
	r.TotalCount = v.TotalCount
	r.Velocity = float32(v.TotalCount / hours)
	r.FirstSeen = v.FirstSeen
	r.LastSeen = v.LastSeen

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
