package repository

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"axon/utils"
	"github.com/lib/pq"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
)

const (
	GetIssueDetailsList               = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, source, destination, COUNT(*) AS total_count, min(time) as first_seen, max(time) as last_seen, ARRAY_AGG(DISTINCT(incident.trace_id)) incidents FROM (select trace_id, source, destination, time from span where workload_id_list is not null and source = ANY($1) and destination = ANY($2)) as s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) GROUP BY issue.issue_hash, issue.issue_title, source, destination, scenario_id, scenario_version LIMIT $3 OFFSET $4"
	GetIssueDetailsByIssueId          = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, source, destination, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY_AGG(DISTINCT(incident.trace_id)) incidents FROM (select * from issue WHERE issue_hash=$1) as issue INNER JOIN incident USING(issue_hash) INNER JOIN (SELECT trace_id, source, destination, time FROM span WHERE workload_id_list IS NOT NULL) AS s USING(trace_id) GROUP BY issue.issue_hash, issue.issue_title, source, destination, scenario_id, scenario_version"
	GetTraceQuery                     = "SELECT trace_id, issue_hash, incident_collection_time from incident where issue_hash=$1 LIMIT $2 OFFSET $3"
	GetSpanRawDataQuery               = "SELECT span.trace_id, span.span_id, request_payload, response_payload, protocol FROM span_raw_data INNER JOIN span USING(span_id) WHERE span.trace_id=$1 AND span.span_id=$2 LIMIT $3 OFFSET $4"
	GetSpanQueryUsingTraceIdAndSpanId = "SELECT trace_id, span_id, source, destination, metadata, latency_ms, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND span_id=$2 LIMIT $3 OFFSET $4"
	GetSpanQueryUsingTraceId          = "SELECT trace_id, span_id, source, destination, metadata, latency_ms, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND workload_id_list is not NULL LIMIT $2 OFFSET $3"
)

var LogTag = "zk_trace_persistence_repo"

type TracePersistenceRepo interface {
	IssueListDetailsRepo(sources, destinations pq.StringArray, offset, limit int) ([]dto.IssueDetailsDto, error)
	GetIssueDetails(issueId string) (dto.IssueDetailsDto, error)
	GetTraces(issueId string, offset, limit int) ([]dto.IncidentTableDto, error)
	GetSpans(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error)
	GetSpanRawData(traceId, spanId string, offset, limit int) ([]dto.SpanRawDataDetailsDto, error)
}

type tracePersistenceRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewTracePersistenceRepo(dbRepo sqlDB.DatabaseRepo) TracePersistenceRepo {
	return &tracePersistenceRepo{dbRepo: dbRepo}
}

func (z tracePersistenceRepo) IssueListDetailsRepo(sources, destinations pq.StringArray, offset, limit int) ([]dto.IssueDetailsDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetIssueDetailsList, []any{sources, destinations, limit, offset})
	defer closeRow()
	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var data []dto.IssueDetailsDto
	for rows.Next() {
		var rawData dto.IssueDetailsDto
		err := rows.Scan(&rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Source, &rawData.Destination, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		days := utils.CalendarDaysBetween(rawData.FirstSeen, rawData.LastSeen) + 1
		rawData.Velocity = float32(rawData.TotalCount / days)
		data = append(data, rawData)
	}

	return data, nil

}

func (z tracePersistenceRepo) GetIssueDetails(issueId string) (dto.IssueDetailsDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetIssueDetailsByIssueId, []any{issueId})
	var res dto.IssueDetailsDto
	defer closeRow()
	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return res, err
	}

	var data []dto.IssueDetailsDto
	for rows.Next() {
		var rawData dto.IssueDetailsDto
		err := rows.Scan(&rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Source, &rawData.Destination, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}

		days := utils.CalendarDaysBetween(rawData.FirstSeen, rawData.LastSeen) + 1
		rawData.Velocity = float32(rawData.TotalCount / days)
		data = append(data, rawData)
	}

	if len(data) > 0 {
		res = data[0]
		zkLogger.Error(LogTag, "more than 1 row returned for issue id ", issueId)
		zkLogger.Error(LogTag, "total rows returned = ", len(data))
	}

	return res, nil
}

func (z tracePersistenceRepo) GetTraces(issueId string, offset, limit int) ([]dto.IncidentTableDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetTraceQuery, []any{issueId, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.IncidentTableDto
	for rows.Next() {
		var rawData dto.IncidentTableDto
		err := rows.Scan(&rawData.TraceId, &rawData.IssueHash, &rawData.IncidentCollectionTime)
		if err != nil {
			zkLogger.Error(LogTag, err)
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
		params = []any{traceId, spanId, limit, offset}
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, params)
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.SpanTableDto
	for rows.Next() {
		var rawData dto.SpanTableDto
		err := rows.Scan(&rawData.TraceId, &rawData.SpanId, &rawData.Source, &rawData.Destination, &rawData.Metadata, &rawData.LatencyMs, &rawData.Protocol, &rawData.Status, &rawData.ParentSpanId, &rawData.WorkloadIdList, &rawData.Time)
		if err != nil {
			zkLogger.Error(LogTag, err)
			return nil, err
		}

		rawData.TraceId = traceId
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetSpanRawData(traceId, spanId string, offset, limit int) ([]dto.SpanRawDataDetailsDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetSpanRawDataQuery, []any{traceId, spanId, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var data []dto.SpanRawDataDetailsDto
	for rows.Next() {
		var rawData dto.SpanRawDataDetailsDto
		err := rows.Scan(&rawData.TraceId, &rawData.SpanId, &rawData.RequestPayload, &rawData.ResponsePayload, &rawData.Protocol)
		if err != nil {
			zkLogger.Error(LogTag, err)
		}
		data = append(data, rawData)
	}

	return data, nil
}
