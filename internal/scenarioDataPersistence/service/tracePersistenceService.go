package service

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/repository"
	"fmt"
	utils "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIssueListWithDetailsService(services string, offset, limit int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetIssueDetailsService(issueHash string) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetIncidentListService(issueHash string, offset, limit int) (traceResponse.IncidentListResponse, *zkErrors.ZkError)
	GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError)
	GetSpanRawDataService(traceId, spanId string) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError)
}

func NewScenarioPersistenceService(repo repository.TracePersistenceRepo) TracePersistenceService {
	return tracePersistenceService{repo: repo}
}

type tracePersistenceService struct {
	repo repository.TracePersistenceRepo
}

func (s tracePersistenceService) GetIssueListWithDetailsService(services string, offset, limit int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueListWithDetailsResponse

	var serviceList []string
	if utils.IsEmpty(services) {
		zkLogger.Info(LogTag, "service list is empty")
	} else {
		l := strings.Split(services, ",")
		for _, service := range l {
			v := strings.TrimSpace(service)
			if utils.IsEmpty(v) {
				continue
			}
			serviceList = append(serviceList, v)
		}
	}

	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, fmt.Sprintf("value of limit or offset is invalid, limit: %d, offset: %d", limit, offset), zkErr)
		return response, &zkErr
	}

	data, err := s.repo.IssueListDetailsRepo(serviceList, offset, limit)
	if err == nil {
		response := traceResponse.ConvertIssueListDetailsDtoToIssueListDetailsResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get issue list with details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIssueDetailsService(issueHash string) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueListWithDetailsResponse

	data, err := s.repo.GetIssueDetails(issueHash)
	if err == nil {
		response := traceResponse.ConvertIssueListDetailsDtoToIssueListDetailsResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get issue details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIncidentListService(issueHash string, offset, limit int) (traceResponse.IncidentListResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentListResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, fmt.Sprintf("value of limit or offset is invalid, limit: %d, offset: %d", limit, offset), zkErr)
		return response, &zkErr
	}

	data, err := s.repo.GetTraces(issueHash, offset, limit)
	if err == nil {
		response := traceResponse.ConvertIncidentTableDtoToIncidentListResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get incident list", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentDetailsResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, fmt.Sprintf("value of limit or offset is invalid, limit: %d, offset: %d", limit, offset), zkErr)
		return response, &zkErr
	}

	data, err := s.repo.GetSpans(traceId, spanId, offset, limit)
	if err == nil {
		response := traceResponse.ConvertSpanToIncidentDetailsResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get incident details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetSpanRawDataService(traceId, spanId string) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError) {
	var response traceResponse.SpanRawDataResponse

	data, err := s.repo.GetSpanRawData(traceId, spanId)
	if err == nil {
		response, respErr := traceResponse.ConvertSpanRawDataToSpanRawDataResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, "failed to convert span raw data to response", err)
			zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
			return *response, &zkErr
		}

		return *response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}
