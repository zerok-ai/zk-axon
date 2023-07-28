package service

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/repository"
	"axon/utils"
	zkErrorsAxon "axon/utils/zkerrors"
	"fmt"
	zkUtils "github.com/zerok-ai/zk-utils-go/common"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
	"time"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIssueListWithDetailsService(services, st string, limit, offset int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetIssueDetailsService(issueHash string) (traceResponse.IssueDetailsResponse, *zkErrors.ZkError)
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

func (s tracePersistenceService) GetIssueListWithDetailsService(services, st string, limit, offset int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueListWithDetailsResponse
	var startTime time.Time
	currentTime := time.Now().UTC()

	if duration, err := utils.ParseTimeString(st); err != nil {
		zkLogger.Error(LogTag, "failed to parse time string", err)
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	} else if currentTime.Add(duration).After(currentTime) {
		zkLogger.Error(LogTag, "time string is not negative", err)
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestStartTimeNotNegative, nil)
		return response, &zkErr
	} else {
		startTime = currentTime.Add(duration)
	}

	var serviceList []string
	if zkUtils.IsEmpty(services) {
		zkLogger.Info(LogTag, "service list is empty")
	} else {
		l := strings.Split(services, ",")
		for _, service := range l {
			v := strings.TrimSpace(service)
			if zkUtils.IsEmpty(v) {
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

	data, err := s.repo.IssueListDetailsRepo(serviceList, startTime, limit, offset)
	if err == nil {
		response := traceResponse.ConvertIssueListDetailsDtoToIssueListDetailsResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get issue list with details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIssueDetailsService(issueHash string) (traceResponse.IssueDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueDetailsResponse

	data, err := s.repo.GetIssueDetails(issueHash)
	if err == nil {
		response := traceResponse.ConvertIssueDetailsDtoToIssueListDetailsResponse(data)
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
