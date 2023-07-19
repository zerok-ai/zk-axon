package repository

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"axon/utils"
	"fmt"
	"github.com/lib/pq"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"strings"
)

const (
	GetIssueDetailsListWithoutServiceName = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, source, destination, COUNT(*) AS total_count, min(time) as first_seen, max(time) as last_seen, ARRAY_AGG( DISTINCT(incident.trace_id) ) incidents FROM ( select trace_id, source, issue_hash_list, destination, time from span where issue_hash_list is not null ) as s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, source, destination, scenario_id, scenario_version LIMIT $1 OFFSET $2"
	GetIssueDetailsList                   = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, source, destination, COUNT(*) AS total_count, min(time) as first_seen, max(time) as last_seen, ARRAY_AGG( DISTINCT(incident.trace_id) ) incidents FROM ( select trace_id, source, issue_hash_list, destination, time from span where issue_hash_list is not null AND (source = ANY($1) OR destination = ANY($2)) ) as s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, source, destination, scenario_id, scenario_version LIMIT $3 OFFSET $4"
	GetIssueDetailsByIssueHash            = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, source, destination, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY_AGG( DISTINCT(incident.trace_id) ) incidents FROM ( select * from issue WHERE issue_hash = $1 ) as issue INNER JOIN incident USING(issue_hash) INNER JOIN ( SELECT trace_id, issue_hash_list, source, destination, time FROM span WHERE issue_hash_list IS NOT NULL ) AS s USING(trace_id) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, source, destination, scenario_id, scenario_version"
	GetTraceQuery                         = "SELECT trace_id, issue_hash, incident_collection_time from incident where issue_hash=$1 LIMIT $2 OFFSET $3"
	GetSpanQueryUsingTraceId              = "SELECT trace_id, span_id, source, destination, metadata, latency_ns, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND workload_id_list is not NULL LIMIT $2 OFFSET $3"
	GetSpanQueryUsingTraceIdAndSpanId     = "SELECT trace_id, span_id, source, destination, metadata, latency_ns, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND span_id=$2"
	GetSpanRawDataQuery                   = "SELECT span.trace_id, span.span_id, request_payload, response_payload, protocol FROM span_raw_data INNER JOIN span USING(span_id) WHERE span.trace_id=$1 AND span.span_id=$2"
)

var LogTag = "zk_trace_persistence_repo"

type TracePersistenceRepo interface {
	IssueListDetailsRepo(serviceList pq.StringArray, offset, limit int) ([]dto.IssueDetailsDto, error)
	GetIssueDetails(issueHash string) ([]dto.IssueDetailsDto, error)
	GetTraces(issueHash string, offset, limit int) ([]dto.IncidentTableDto, error)
	GetSpans(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error)
	GetSpanRawData(traceId, spanId string) ([]dto.SpanRawDataDetailsDto, error)
}

type tracePersistenceRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewTracePersistenceRepo(dbRepo sqlDB.DatabaseRepo) TracePersistenceRepo {
	return &tracePersistenceRepo{dbRepo: dbRepo}
}

func (z tracePersistenceRepo) IssueListDetailsRepo(serviceList pq.StringArray, offset, limit int) ([]dto.IssueDetailsDto, error) {
	var query string
	var params []any
	if serviceList == nil || len(serviceList) == 0 {
		query = GetIssueDetailsListWithoutServiceName
		params = []any{limit, offset}
	} else {
		query = GetIssueDetailsList
		params = []any{serviceList, serviceList, limit, offset}
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, params)
	defer closeRow()
	if err != nil || rows == nil {
		s := strings.Join(serviceList, ",")
		zkLogger.Error(LogTag, fmt.Sprintf("service_list: %s", s), err)
		return nil, err
	}

	data := make([]dto.IssueDetailsDto, 0)
	for rows.Next() {
		var rawData dto.IssueDetailsDto
		err := rows.Scan(&rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Source, &rawData.Destination, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			s := strings.Join(serviceList, ",")
			zkLogger.Error(LogTag, fmt.Sprintf("service_list: %s", s), err)
		}

		hours := utils.HoursBetween(rawData.FirstSeen, rawData.LastSeen) + 1
		rawData.Velocity = float32(rawData.TotalCount / hours)
		data = append(data, rawData)
	}

	return data, nil

}

func (z tracePersistenceRepo) GetIssueDetails(issueHash string) ([]dto.IssueDetailsDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetIssueDetailsByIssueHash, []any{issueHash})
	defer closeRow()
	if err != nil || rows == nil {
		zkLogger.Error(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), err)
		return nil, err
	}

	data := make([]dto.IssueDetailsDto, 0)
	for rows.Next() {
		var rawData dto.IssueDetailsDto
		err := rows.Scan(&rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Source, &rawData.Destination, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			zkLogger.Error(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), err)
		}

		hours := utils.HoursBetween(rawData.FirstSeen, rawData.LastSeen) + 1
		rawData.Velocity = float32(rawData.TotalCount / hours)
		data = append(data, rawData)
	}

	return data, nil
}

func (z tracePersistenceRepo) GetTraces(issueHash string, offset, limit int) ([]dto.IncidentTableDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetTraceQuery, []any{issueHash, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), err)
		return nil, err
	}

	var responseArr []dto.IncidentTableDto
	for rows.Next() {
		var rawData dto.IncidentTableDto
		err := rows.Scan(&rawData.TraceId, &rawData.IssueHash, &rawData.IncidentCollectionTime)
		if err != nil {
			zkLogger.Error(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), err)
		}

		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetSpans(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error) {
	var query string
	var params []any
	if zkCommon.IsEmpty(spanId) {
		query = GetSpanQueryUsingTraceId
		params = []any{traceId, limit, offset}
	} else {
		query = GetSpanQueryUsingTraceIdAndSpanId
		params = []any{traceId, spanId}
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, params)
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, fmt.Sprintf("trace_id: %s, span_id: %s", traceId, spanId), err)
		return nil, err
	}

	var responseArr []dto.SpanTableDto
	for rows.Next() {
		var rawData dto.SpanTableDto
		err := rows.Scan(&rawData.TraceId, &rawData.SpanId, &rawData.Source, &rawData.Destination, &rawData.Metadata, &rawData.LatencyNs, &rawData.Protocol, &rawData.Status, &rawData.ParentSpanId, &rawData.WorkloadIdList, &rawData.Time)
		if err != nil {
			zkLogger.Error(LogTag, fmt.Sprintf("trace_id: %s, span_id: %s", traceId, spanId), err)
			return nil, err
		}

		rawData.TraceId = traceId
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetSpanRawData(traceId, spanId string) ([]dto.SpanRawDataDetailsDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetSpanRawDataQuery, []any{traceId, spanId})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, fmt.Sprintf("trace_id: %s, span_id: %s", traceId, spanId), err)
		return nil, err
	}

	var data []dto.SpanRawDataDetailsDto
	for rows.Next() {
		var rawData dto.SpanRawDataDetailsDto
		err := rows.Scan(&rawData.TraceId, &rawData.SpanId, &rawData.RequestPayload, &rawData.ResponsePayload, &rawData.Protocol)
		if err != nil {
			zkLogger.Error(LogTag, fmt.Sprintf("trace_id: %s, span_id: %s", traceId, spanId), err)
		}
		data = append(data, rawData)
	}

	return data, nil
}
