package service

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/repository"
	"axon/utils"
	"fmt"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIncidentData(scenarioType, source string, offset, limit int) (traceResponse.IncidentResponse, *zkErrors.ZkError)
	GetTraces(scenarioId string, offset, limit int) (traceResponse.TraceResponse, *zkErrors.ZkError)
	GetTracesMetadata(traceId, spanId string, offset, limit int) (traceResponse.SpanResponse, *zkErrors.ZkError)
	GetTracesRawData(traceId, spanId string, offset, limit int) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError)
	GetMetadataMap(duration string, offset, limit int) (traceResponse.MetadataMapResponse, *zkErrors.ZkError)
}

func NewScenarioPersistenceService(repo repository.TracePersistenceRepo) TracePersistenceService {
	return tracePersistenceService{repo: repo}
}

type tracePersistenceService struct {
	repo repository.TracePersistenceRepo
}

func (s tracePersistenceService) GetIncidentData(scenarioType, source string, offset, limit int) (traceResponse.IncidentResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetIncidentData(scenarioType, source, offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertIncidentToIncidentResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetTraces(scenarioId string, offset, limit int) (traceResponse.TraceResponse, *zkErrors.ZkError) {
	var response traceResponse.TraceResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetTraces(scenarioId, offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertScenarioTableDtoToTraceResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil
	}
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetTracesMetadata(traceId, spanId string, offset, limit int) (traceResponse.SpanResponse, *zkErrors.ZkError) {
	var response traceResponse.SpanResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetSpan(traceId, spanId, offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertSpanToSpanResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetTracesRawData(traceId, spanId string, offset, limit int) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError) {
	var response traceResponse.SpanRawDataResponse
	//TODO: discuss if the below condition of limit > 100 is fine. or it should be read from some config
	threshold := 100
	if offset < 0 || limit < 1 || limit > threshold {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, fmt.Sprintf("either offset or limit < 0 or limit > %d", threshold))
		return response, &zkErr
	}

	data, err := s.repo.GetSpanRawData(traceId, spanId, offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertSpanRawDataToSpanRawDataResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil

	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetMetadataMap(duration string, offset, limit int) (traceResponse.MetadataMapResponse, *zkErrors.ZkError) {
	var response traceResponse.MetadataMapResponse
	if !utils.IsValidPxlTime(duration) {
		return response, zkCommon.ToPtr(zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, "invalid duration"))
	}

	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetMetadataMap("st", offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertMetadataMapToMetadataMapResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil
	}
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr

}
