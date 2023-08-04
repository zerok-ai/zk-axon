package repository

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"fmt"
	"github.com/lib/pq"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"strings"
	"time"
)

const (
	GetIssueDetailsListWithoutServiceNameAndScenarioIdFilter     = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, issue.issue_hash, issue.issue_title, scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY (SELECT temp.trace_id FROM (SELECT DISTINCT i1.trace_id, max(s1.time) FROM incident i1 INNER JOIN span s1 using(trace_id) WHERE i1.trace_id=ANY(ARRAY_AGG(incident.trace_id)) GROUP BY i1.trace_id ORDER BY max(s1.time) DESC) AS TEMP) incidents FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $2) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, scenario_id, scenario_version ORDER BY last_seen DESC LIMIT $3 OFFSET $4"
	GetIssueDetailsListWithoutServiceNameAndWithScenarioIdFilter = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, issue.issue_hash, issue.issue_title, scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY (SELECT temp.trace_id FROM (SELECT DISTINCT i1.trace_id, max(s1.time) FROM incident i1 INNER JOIN span s1 using(trace_id) WHERE i1.trace_id=ANY(ARRAY_AGG(incident.trace_id)) GROUP BY i1.trace_id ORDER BY max(s1.time) DESC) AS TEMP) incidents FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $2) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) AND scenario_id=ANY($3) GROUP BY issue.issue_hash, issue.issue_title, scenario_id, scenario_version ORDER BY last_seen DESC LIMIT $4 OFFSET $5"
	GetIssueDetailsWithServiceNameAndWithoutScenarioIdListFilter = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, issue.issue_hash, issue.issue_title, scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY (SELECT temp.trace_id FROM (SELECT DISTINCT i1.trace_id, max(s1.time) FROM incident i1 INNER JOIN span s1 using(trace_id) WHERE i1.trace_id=ANY(ARRAY_AGG(incident.trace_id)) GROUP BY i1.trace_id ORDER BY max(s1.time) DESC) AS TEMP) incidents FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $2 AND (SOURCE = ANY($3) OR destination = ANY($4))) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, scenario_id, scenario_version ORDER BY last_seen DESC LIMIT $5 OFFSET $6"
	GetIssueDetailsListWithServiceNameAndScenarioIdFilter        = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, issue.issue_hash, issue.issue_title, scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY (SELECT temp.trace_id FROM (SELECT DISTINCT i1.trace_id, max(s1.time) FROM incident i1 INNER JOIN span s1 using(trace_id) WHERE i1.trace_id=ANY(ARRAY_AGG(incident.trace_id)) GROUP BY i1.trace_id ORDER BY max(s1.time) DESC) AS TEMP) incidents FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $2 AND (SOURCE = ANY($3) OR destination = ANY($4))) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) AND scenario_id=ANY($5) GROUP BY issue.issue_hash, issue.issue_title, scenario_id, scenario_version ORDER BY last_seen DESC LIMIT $6 OFFSET $7"
	GetScenarioDetailsWithoutServiceNameFilter                   = "SELECT scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $1) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) AND scenario_id=ANY($2) GROUP BY scenario_id, scenario_version ORDER BY last_seen DESC"
	GetScenarioDetailsWithServiceNameFilter                      = "SELECT scenario_id, scenario_version, ARRAY_AGG(DISTINCT(SOURCE)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen FROM (SELECT trace_id, SOURCE, issue_hash_list, destination, time FROM span WHERE issue_hash_list IS NOT NULL AND time > $1 AND (SOURCE = ANY($2) OR destination = ANY($3))) AS s INNER JOIN incident USING(trace_id) INNER JOIN issue USING(issue_hash) WHERE issue.issue_hash = ANY(issue_hash_list) AND scenario_id=ANY($4) GROUP BY scenario_id, scenario_version ORDER BY last_seen DESC"
	GetIssueDetailsByIssueHash                                   = "SELECT issue.issue_hash, issue.issue_title, scenario_id, scenario_version, ARRAY_AGG(DISTINCT(source)) sources, ARRAY_AGG(DISTINCT(destination)) destinations, COUNT(*) AS total_count, min(time) AS first_seen, max(time) AS last_seen, ARRAY (SELECT temp.trace_id FROM (SELECT DISTINCT i1.trace_id, max(s1.time) FROM incident i1 INNER JOIN span s1 using(trace_id) WHERE i1.trace_id=ANY(ARRAY_AGG(incident.trace_id)) GROUP BY i1.trace_id ORDER BY max(s1.time) DESC)  AS TEMP) incidents FROM (SELECT * FROM issue WHERE issue_hash = $1) AS issue INNER JOIN incident USING(issue_hash) INNER JOIN (SELECT trace_id, issue_hash_list, source, destination, time FROM span WHERE issue_hash_list IS NOT NULL) AS s USING(trace_id) WHERE issue.issue_hash = ANY(issue_hash_list) GROUP BY issue.issue_hash, issue.issue_title, scenario_id, scenario_version"
	GetTraceQuery                                                = "SELECT CASE WHEN $1 then COUNT(*) OVER() ELSE 0 END AS total_rows, trace_id, issue_hash, incident_collection_time, entry_service, end_point, protocol, root_span_time, latency_ns from incident where issue_hash=$2 ORDER BY incident_collection_time DESC LIMIT $3 OFFSET $4"
	GetSpanQueryUsingTraceId                                     = "SELECT trace_id, span_id, source, destination, metadata, latency_ns, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND workload_id_list is not NULL ORDER BY time DESC LIMIT $2 OFFSET $3"
	GetSpanQueryUsingTraceIdAndSpanId                            = "SELECT trace_id, span_id, source, destination, metadata, latency_ns, protocol, status, parent_span_id, workload_id_list, time FROM span WHERE trace_id=$1 AND span_id=$2"
	GetSpanRawDataQuery                                          = "SELECT span.trace_id, span.span_id, request_payload, response_payload, protocol FROM span_raw_data INNER JOIN span USING(span_id) WHERE span.trace_id=$1 AND span.span_id=$2"
	GetTraceQueryByScenarioId                                    = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, trace_id, incident_collection_time, entry_service, end_point, protocol, root_span_time, latency_ns FROM (SELECT DISTINCT ON (incident.trace_id) incident.trace_id, incident.incident_collection_time, incident.entry_service, incident.end_point, incident.protocol, incident.root_span_time, incident.latency_ns FROM issue INNER JOIN incident ON issue.issue_hash = incident.issue_hash WHERE issue.scenario_id = $2 AND issue.issue_hash=$3) AS distinct_incidents ORDER BY incident_collection_time DESC LIMIT $4 OFFSET $5"
	GetTraceQueryByScenarioIdWithoutIssueHashFilter              = "SELECT CASE WHEN $1 THEN COUNT(*) OVER() ELSE 0 END AS total_rows, trace_id, incident_collection_time, entry_service, end_point, protocol, root_span_time, latency_ns FROM (SELECT DISTINCT ON (incident.trace_id) incident.trace_id, incident.incident_collection_time, incident.entry_service, incident.end_point, incident.protocol, incident.root_span_time, incident.latency_ns FROM issue INNER JOIN incident ON issue.issue_hash = incident.issue_hash WHERE issue.scenario_id = $2) AS distinct_incidents ORDER BY incident_collection_time DESC LIMIT $3 OFFSET $4"
)

var LogTag = "zk_trace_persistence_repo"

type TracePersistenceRepo interface {
	IssueListDetailsRepo(serviceList pq.StringArray, scenarioList []int, limit, offset int, st time.Time) ([]dto.IssueDetailsDto, error)
	GetScenarioDetailsRepo(scenarioId, serviceList pq.StringArray, st time.Time) ([]dto.ScenarioDetailsDto, error)
	GetIssueDetails(issueHash string) ([]dto.IssueDetailsDto, error)
	GetTraces(issueHash string, offset, limit int) ([]dto.IncidentTableDto, error)
	GetTracesForScenarioId(scenarioId, issueHash string, limit, offset int) ([]dto.IncidentTableDto, error)
	GetSpans(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error)
	GetSpanRawData(traceId, spanId string) ([]dto.SpanRawDataDetailsDto, error)
}

type tracePersistenceRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func (z tracePersistenceRepo) GetTracesForScenarioId(scenarioId, issueHash string, limit, offset int) ([]dto.IncidentTableDto, error) {
	var query string
	var params []any
	if zkCommon.IsEmpty(issueHash) {
		query = GetTraceQueryByScenarioIdWithoutIssueHashFilter
		params = []any{true, scenarioId, limit, offset}
	} else {
		query = GetTraceQueryByScenarioId
		params = []any{true, scenarioId, issueHash, limit, offset}
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, params)

	if rows != nil {
		defer closeRow()
	}

	logMessage := fmt.Sprintf("scenario Id: %s", scenarioId)

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, logMessage, err)
		return nil, err
	}

	var responseArr []dto.IncidentTableDto
	for rows.Next() {
		var rawData dto.IncidentTableDto
		err := rows.Scan(&rawData.TotalRows, &rawData.TraceId, &rawData.IncidentCollectionTime, &rawData.EntryService, &rawData.EndPoint, &rawData.Protocol, &rawData.RootSpanTime, &rawData.LatencyNs)
		if err != nil {
			zkLogger.Error(LogTag, logMessage, err)
		}
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func NewTracePersistenceRepo(dbRepo sqlDB.DatabaseRepo) TracePersistenceRepo {
	return &tracePersistenceRepo{dbRepo: dbRepo}
}

func (z tracePersistenceRepo) IssueListDetailsRepo(serviceList pq.StringArray, scenarioList []int, limit, offset int, st time.Time) ([]dto.IssueDetailsDto, error) {
	var query string
	var params []any

	var scenarioListInt32 pq.Int32Array
	for _, i := range scenarioList {
		scenarioListInt32 = append(scenarioListInt32, int32(i))
	}

	if serviceList == nil || len(serviceList) == 0 {
		if scenarioListInt32 == nil || len(scenarioListInt32) == 0 {
			query = GetIssueDetailsListWithoutServiceNameAndScenarioIdFilter
			params = []any{true, st, limit, offset}
		} else {
			query = GetIssueDetailsListWithoutServiceNameAndWithScenarioIdFilter
			params = []any{true, st, scenarioListInt32, limit, offset}
		}
	} else {
		if scenarioListInt32 == nil || len(scenarioListInt32) == 0 {
			query = GetIssueDetailsWithServiceNameAndWithoutScenarioIdListFilter
			params = []any{true, st, serviceList, serviceList, limit, offset}
		} else {
			query = GetIssueDetailsListWithServiceNameAndScenarioIdFilter
			params = []any{true, st, serviceList, serviceList, scenarioListInt32, limit, offset}
		}
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
		err := rows.Scan(&rawData.TotalRows, &rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Sources, &rawData.Destinations, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			s := strings.Join(serviceList, ",")
			zkLogger.Error(LogTag, fmt.Sprintf("service_list: %s", s), err)
		}

		data = append(data, rawData)
	}

	return data, nil

}

func (z tracePersistenceRepo) GetScenarioDetailsRepo(scenarioId, serviceList pq.StringArray, st time.Time) ([]dto.ScenarioDetailsDto, error) {
	var query string
	var params []any
	if serviceList == nil || len(serviceList) == 0 {
		query = GetScenarioDetailsWithoutServiceNameFilter
		params = []any{st, scenarioId}
	} else {
		query = GetScenarioDetailsWithServiceNameFilter
		params = []any{st, serviceList, serviceList, scenarioId}
	}

	rows, err, closeRow := z.dbRepo.GetAll(query, params)
	defer closeRow()
	if err != nil || rows == nil {
		s := strings.Join(serviceList, ",")
		sc := strings.Join(scenarioId, ",")
		zkLogger.Error(LogTag, fmt.Sprintf("service_list: %s, scenario_id_list: %s", s, sc), err)
		return nil, err
	}

	data := make([]dto.ScenarioDetailsDto, 0)
	for rows.Next() {
		var rawData dto.ScenarioDetailsDto
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Sources, &rawData.Destinations, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen)
		if err != nil {
			s := strings.Join(serviceList, ",")
			sc := strings.Join(scenarioId, ",")
			zkLogger.Error(LogTag, fmt.Sprintf("service_list: %s, scenario_id_list: %s", s, sc), err)
		}

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
		err := rows.Scan(&rawData.IssueHash, &rawData.IssueTitle, &rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Sources, &rawData.Destinations, &rawData.TotalCount, &rawData.FirstSeen, &rawData.LastSeen, &rawData.Incidents)
		if err != nil {
			zkLogger.Error(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), err)
		}

		data = append(data, rawData)
	}

	if len(data) > 1 {
		zkLogger.Info(LogTag, fmt.Sprintf("issue_hash: %s", issueHash), "multiple rows found for issue hash")
	}

	return data, nil
}

func (z tracePersistenceRepo) GetTraces(issueHash string, offset, limit int) ([]dto.IncidentTableDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetTraceQuery, []any{true, issueHash, limit, offset})
	defer closeRow()

	logMessage := fmt.Sprintf("issue_hash: %s", issueHash)

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, logMessage, err)
		return nil, err
	}

	var responseArr []dto.IncidentTableDto
	for rows.Next() {
		var rawData dto.IncidentTableDto
		err := rows.Scan(&rawData.TotalRows, &rawData.TraceId, &rawData.IssueHash, &rawData.IncidentCollectionTime, &rawData.EntryService, &rawData.EndPoint, &rawData.Protocol, &rawData.RootSpanTime, &rawData.LatencyNs)
		if err != nil {
			zkLogger.Error(LogTag, logMessage, err)
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
