package service

import (
	traceResponse "axon/internal/scenarioDataPersistence/model/response"
	"axon/internal/scenarioDataPersistence/repository"
	"fmt"
	"github.com/lib/pq"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIssueListWithDetailsService(source, destination string, offset, limit int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetIssueDetailsService(issueId string) (traceResponse.IssueWithDetailsResponse, *zkErrors.ZkError)
	GetIncidentListService(issueId string, offset, limit int) (traceResponse.IncidentListResponse, *zkErrors.ZkError)
	GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError)
	GetSpanRawDataService(traceId, spanId string, offset, limit int) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError)
}

func NewScenarioPersistenceService(repo repository.TracePersistenceRepo) TracePersistenceService {
	return tracePersistenceService{repo: repo}
}

type tracePersistenceService struct {
	repo repository.TracePersistenceRepo
}

func (s tracePersistenceService) GetIssueListWithDetailsService(sources, destinations string, offset, limit int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueListWithDetailsResponse

	sourceList := strings.Split(sources, ",")
	if sourceList == nil || len(sourceList) == 0 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, "source is empty")
		zkLogger.Error(LogTag, zkErr, sourceList)
		return response, &zkErr
	}

	destinationList := strings.Split(destinations, ",")
	if destinationList == nil || len(destinationList) == 0 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, "destination is empty")
		zkLogger.Error(LogTag, zkErr, destinationList)
		return response, &zkErr
	}

	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return response, &zkErr
	}

	x := pq.StringArray(sourceList)
	y := pq.StringArray(destinationList)

	data, err := s.repo.IssueListDetailsRepo(x, y, offset, limit)
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
