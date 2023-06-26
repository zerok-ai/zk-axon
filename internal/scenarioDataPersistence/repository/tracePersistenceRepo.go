package repository

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"axon/utils"
	"log"
	"time"

	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
)

const (
	ScenarioTraceTablePostgres = "scenario_trace"
	SpanTablePostgres          = "span"
	SpanRawDataTablePostgres   = "span_raw_data"
	ErrorFlag                  = true

	ScenarioId      = "scenario_id"
	ScenarioVersion = "scenario_version"
	TraceId         = "trace_id"
	ScenarioTitle   = "scenario_title"
	ScenarioType    = "scenario_type"

	SpanId          = "span_id"
	ParentSpanId    = "parent_span_id"
	Source          = "source"
	Destination     = "destination"
	WorkloadIdList  = "workload_id_list"
	Metadata        = "metadata"
	LatencyMs       = "latency_ms"
	Protocol        = "protocol"
	RequestPayload  = "request_payload"
	ResponsePayload = "response_payload"

	GetIncidentData                   = "SELECT t.scenario_id, t.scenario_version, t.scenario_title, COUNT(*) AS incident_count, md.destination, min(t.created_at) as first_seen, max(t.created_at) as last_seen FROM (SELECT * FROM scenario_trace WHERE scenario_type=$1) AS t INNER JOIN (SELECT * FROM span WHERE source=$2) AS md USING(trace_id) GROUP BY t.scenario_id, t.scenario_version, t.scenario_title, md.destination  LIMIT $3 OFFSET $4"
	GetTraceQuery                     = "SELECT scenario_version, trace_id FROM scenario_trace WHERE scenario_id=$1 LIMIT $2 OFFSET $3"
	GetSpanRawDataQuery               = "SELECT request_payload, response_payload FROM span_raw_data WHERE trace_id=$1 AND span_id=$2 LIMIT $3 OFFSET $4"
	GetSpanQueryUsingTraceIdAndSpanId = "SELECT span_id, parent_span_id, source, destination, workload_id_list, metadata, latency_ms, protocol FROM span WHERE trace_id=$1 AND span_id=$2 LIMIT $3 OFFSET $4"
	GetSpanQueryUsingTraceId          = "SELECT span_id, parent_span_id, source, destination, workload_id_list, metadata, latency_ms, protocol FROM span WHERE trace_id=$1 LIMIT $2 OFFSET $3"
	GetMetadataMapQueryUsingDuration  = "SELECT md.source, md.destination, COUNT(DISTINCT(trace_id)) AS trace_count,  ARRAY_AGG(DISTINCT(protocol)) protocol_list FROM (SELECT trace_id AS trace_id FROM scenario_trace WHERE created_at >= $1) AS st INNER JOIN (SELECT trace_id, protocol, source, destination FROM span WHERE workload_id_list IS NOT NULL) AS md USING(trace_id) GROUP BY md.source, md.destination LIMIT $2 OFFSET $3"
)

var LogTag = "zk_trace_persistence_repo"

type TracePersistenceRepo interface {
	GetIncidentData(errorType, source string, offset, limit int) ([]dto.IncidentDto, error)
	GetTraces(scenarioId string, offset, limit int) ([]dto.ScenarioTableDto, error)
	GetSpan(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error)
	GetSpanRawData(traceId, spanId string, offset, limit int) ([]dto.SpanRawDataTableDto, error)
	GetMetadataMap(st string, offset, limit int) ([]dto.MetadataMapDto, error)
}

type tracePersistenceRepo struct {
	dbRepo sqlDB.DatabaseRepo
}

func NewTracePersistenceRepo(dbRepo sqlDB.DatabaseRepo) TracePersistenceRepo {
	return &tracePersistenceRepo{dbRepo: dbRepo}
}

func (z tracePersistenceRepo) GetIncidentData(scenarioType, source string, offset, limit int) ([]dto.IncidentDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetIncidentData, []any{scenarioType, source, limit, offset})
	defer closeRow()
	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.IncidentDto
	for rows.Next() {
		var rawData dto.IncidentDto
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.Title, &rawData.TotalCount,
			&rawData.Destination, &rawData.FirstSeen, &rawData.LastSeen)
		if err != nil {
			log.Fatal(err)
		}

		rawData.ScenarioType = scenarioType
		rawData.Source = source

		last, err := utils.ParseTimestamp(rawData.LastSeen)
		if err != nil {
			zkLogger.Error(LogTag, "unable to parse last seen:", rawData.LastSeen, ", err:", err)
			rawData.Velocity = -1
			continue
		}

		first, err := utils.ParseTimestamp(rawData.FirstSeen)
		if err != nil {
			zkLogger.Error(LogTag, "unable to parse first seen:", rawData.FirstSeen, ", err:", err)
			rawData.Velocity = -1
			continue
		}

		days := utils.CalendarDaysBetween(first, last) + 1
		rawData.Velocity = float32(rawData.TotalCount / days)
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil

}

func (z tracePersistenceRepo) GetTraces(scenarioId string, offset, limit int) ([]dto.ScenarioTableDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetTraceQuery, []any{scenarioId, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.ScenarioTableDto
	for rows.Next() {
		var rawData dto.ScenarioTableDto
		err := rows.Scan(&rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}
		rawData.ScenarioId = scenarioId

		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetSpan(traceId, spanId string, offset, limit int) ([]dto.SpanTableDto, error) {
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
		err := rows.Scan(&rawData.SpanId, &rawData.ParentSpanId, &rawData.Source, &rawData.Destination, &rawData.WorkloadIdList, &rawData.Metadata, &rawData.LatencyMs, &rawData.Protocol)
		if err != nil {
			log.Fatal(err)
		}

		rawData.TraceId = traceId
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetSpanRawData(traceId, spanId string, offset, limit int) ([]dto.SpanRawDataTableDto, error) {
	rows, err, closeRow := z.dbRepo.GetAll(GetSpanRawDataQuery, []any{traceId, spanId, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.SpanRawDataTableDto
	for rows.Next() {
		var rawData dto.SpanRawDataTableDto
		err := rows.Scan(&rawData.RequestPayload, &rawData.ResponsePayload)
		if err != nil {
			log.Fatal(err)
		}
		rawData.TraceId = traceId
		rawData.SpanId = spanId
		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}

func (z tracePersistenceRepo) GetMetadataMap(st string, offset, limit int) ([]dto.MetadataMapDto, error) {
	twoMinutesAgo := time.Now().Add(-40000 * time.Minute)
	rows, err, closeRow := z.dbRepo.GetAll(GetMetadataMapQueryUsingDuration, []any{twoMinutesAgo, limit, offset})
	defer closeRow()

	if err != nil || rows == nil {
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	var responseArr []dto.MetadataMapDto
	for rows.Next() {
		var rawData dto.MetadataMapDto
		err := rows.Scan(&rawData.Source, &rawData.Destination, &rawData.TraceCount, &rawData.ProtocolList)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	return responseArr, nil
}
