package service

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/repository"
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIssueListWithDetailsService(services string, offset, limit int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetIssueDetailsService(issueId string) (traceResponse.IssueWithDetailsResponse, *zkErrors.ZkError)
	GetIncidentListService(issueId string, offset, limit int) (traceResponse.IncidentListResponse, *zkErrors.ZkError)
	GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError)
	GetSpanRawDataService(traceId, spanId string, offset, limit int) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError)
	//GetScenariosAllTracesDataService(scenarioId string, offset, limit int) (traceResponse.ScenarioIncidentDetailsResponse, *zkErrors.ZkError)
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
	if serviceList == nil || len(serviceList) == 0 {
		zkLogger.Info(LogTag, "service list is empty")
	} else {
		serviceList = strings.Split(services, ",")
	}

	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.IssueListDetailsRepo(serviceList, offset, limit)
	if err == nil {
		response := traceResponse.ConvertIssueListDetailsDtoToIssueListDetailsResponse(data)
		return *response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIssueDetailsService(issueId string) (traceResponse.IssueWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueWithDetailsResponse

	data, err := s.repo.GetIssueDetails(issueId)
	if err == nil {
		response := traceResponse.ConvertIssueToIssueDetailsResponse(data)
		return response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIncidentListService(issueId string, offset, limit int) (traceResponse.IncidentListResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentListResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetTraces(issueId, offset, limit)
	if err == nil {
		response := traceResponse.ConvertIncidentTableDtoToIncidentListResponse(data)
		return *response, nil
	}
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentDetailsResponse
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetSpans(traceId, spanId, offset, limit)
	if err == nil {
		response, respErr := traceResponse.ConvertSpanToIncidentDetailsResponse(data)
		if respErr != nil {
			zkLogger.Error(LogTag, err)
		}
		return *response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetSpanRawDataService(traceId, spanId string, offset, limit int) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError) {
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

//func (s tracePersistenceService) GetScenariosAllTracesDataService(scenarioId string, offset, limit int) (traceResponse.ScenarioIncidentDetailsResponse, *zkErrors.ZkError) {
//	var response traceResponse.ScenarioIncidentDetailsResponse
//	//TODO: discuss if the below condition of limit > 100 is fine. or it should be read from some config
//	threshold := 1000
//	if offset < 0 || limit < 1 || limit > threshold {
//		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, fmt.Sprintf("either offset or limit < 0 or limit > %d", threshold))
//		return response, &zkErr
//	}
//
//	data, err := s.repo.GetScenariosAllTracesDataService(scenarioId, offset, limit)
//	if err == nil {
//		response := traceResponse.ConvertScenarioIncidentDetailsDtoListToScenarioIncidentDetailsResponse(data)
//		return response, nil
//	}
//
//	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
//	return response, &zkErr
//}
