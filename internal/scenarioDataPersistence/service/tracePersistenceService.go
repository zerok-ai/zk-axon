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
	"strconv"
	"strings"
	"time"
)

var LogTag = "zk_trace_persistence_service"

type TracePersistenceService interface {
	GetIssueListWithDetailsService(services, scenarioIds, st string, limit, offset int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError)
	GetScenarioDetailsService(scenarioIds, services, st string) (traceResponse.ScenarioDetailsResponse, *zkErrors.ZkError)
	GetIssueDetailsService(issueHash string) (traceResponse.IssueDetailsResponse, *zkErrors.ZkError)
	GetIncidentListService(issueHash string, offset, limit int) (traceResponse.IncidentIdListResponse, *zkErrors.ZkError)
	GetIncidentDetailsService(traceId, spanId string, offset, limit int) (traceResponse.IncidentDetailsResponse, *zkErrors.ZkError)
	GetSpanRawDataService(traceId, spanId string) (traceResponse.SpanRawDataResponse, *zkErrors.ZkError)
	GetIncidentListServiceForScenarioId(scenarioId, issueHash string, offset, limit int) (traceResponse.IncidentDetailListResponse, *zkErrors.ZkError)
}

func NewScenarioPersistenceService(repo repository.TracePersistenceRepo) TracePersistenceService {
	return tracePersistenceService{repo: repo}
}

type tracePersistenceService struct {
	repo repository.TracePersistenceRepo
}

func (s tracePersistenceService) GetIncidentListServiceForScenarioId(scenarioId, issueHash string, offset, limit int) (traceResponse.IncidentDetailListResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentDetailListResponse

	if zkErr := utils.ValidateOffsetLimitValue(offset, limit); zkErr != nil {
		return response, zkErr
	}

	data, err := s.repo.GetTracesForScenarioId(scenarioId, issueHash, limit, offset)
	if err == nil {
		response := traceResponse.ConvertIncidentTableDtoToIncidentDetailListResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get incident list for scenario Id", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIssueListWithDetailsService(services, scenarioIds, st string, limit, offset int) (traceResponse.IssueListWithDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.IssueListWithDetailsResponse
	var startTime time.Time

	//R: We should do this validation in the handler method along with other validations, and we can pass the duration to this method.
	// My idea was to only validation the syntactic parts in handler and actual business logic in service layer.

	// My idea was to only validation the syntactic parts in handler and actual business logic in service layer.
	// Done
	if t, zkErr := getStartTime(st); zkErr != nil {
		return response, zkErr
	} else {
		startTime = t
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

	var scenarioIdList []int
	if zkUtils.IsEmpty(scenarioIds) {
		zkLogger.Info(LogTag, "scenarioIds list is empty")
	} else {
		l := strings.Split(scenarioIds, ",")
		for _, scenario := range l {
			v := strings.TrimSpace(scenario)
			if zkUtils.IsEmpty(v) {
				continue
			}
			i, err := strconv.Atoi(v)
			if err != nil {
				zkLogger.Error(LogTag, fmt.Sprintf("failed to convert scenario id to int, scenarioId: %s ", v), err)
				zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestScenarioIdNotInteger, nil)
				return response, &zkErr
			}
			scenarioIdList = append(scenarioIdList, i)
		}
	}

	//R: This looks like generic validation. Can we move this to a utils method?
	//Done
	if zkErr := utils.ValidateOffsetLimitValue(offset, limit); zkErr != nil {
		return response, zkErr
	}

	data, err := s.repo.IssueListDetailsRepo(serviceList, scenarioIdList, limit, offset, startTime)
	if err == nil {
		response := traceResponse.ConvertIssueListDetailsDtoToIssueListDetailsResponse(data)
		return *response, nil
	}

	zkLogger.Error(LogTag, "failed to get issue list with details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func getStartTime(st string) (time.Time, *zkErrors.ZkError) {
	var t time.Time
	currentTime := time.Now().UTC()

	if duration, err := utils.ParseTimeString(st); err != nil {
		zkLogger.Error(LogTag, "failed to parse time string", err)
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		return t, &zkErr
	} else if currentTime.Add(duration).After(currentTime) {
		zkLogger.Error(LogTag, "time string is not negative", err)
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestStartTimeNotNegative, nil)
		return t, &zkErr
	} else {
		t = currentTime.Add(duration)
	}

	return t, nil
}

func (s tracePersistenceService) GetScenarioDetailsService(scenarioIds, services, st string) (traceResponse.ScenarioDetailsResponse, *zkErrors.ZkError) {
	var response traceResponse.ScenarioDetailsResponse
	var startTime time.Time

	//R: We should do this validation in the handler method along with other validations, and we can pass the duration to this method.
	// My idea was to only validation the syntactic parts in handler and actual business logic in service layer.

	//R: The same code is repeated above in GetIssueListWithDetailsService, we can move it to a common method.
	// Done
	if t, zkErr := getStartTime(st); zkErr != nil {
		return response, zkErr
	} else {
		startTime = t
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

	scenarioIdList := make([]string, 0)
	if zkUtils.IsEmpty(scenarioIds) {
		zkLogger.Error(LogTag, "scenario id list is empty")
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestScenarioIdListEmpty, nil)
		return response, &zkErr
	} else {
		l := strings.Split(scenarioIds, ",")
		for _, scenario := range l {
			v := strings.TrimSpace(scenario)
			if zkUtils.IsEmpty(v) {
				continue
			}
			scenarioIdList = append(scenarioIdList, v)
		}
	}

	//R: nil check is not required. Since len will be 0 even for nil array.
	// Done
	if len(scenarioIdList) == 0 {
		zkLogger.Error(LogTag, "scenario id list is empty after parsing")
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestScenarioIdListEmpty, nil)
		return response, &zkErr
	}

	data, err := s.repo.GetScenarioDetailsRepo(scenarioIdList, serviceList, startTime)
	if err == nil {
		response := traceResponse.ConvertScenarioDetailsDtoToScenarioDetailsResponse(data)
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
		return response, nil
	}

	zkLogger.Error(LogTag, "failed to get issue details", err)
	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}

func (s tracePersistenceService) GetIncidentListService(issueHash string, offset, limit int) (traceResponse.IncidentIdListResponse, *zkErrors.ZkError) {
	var response traceResponse.IncidentIdListResponse
	//R: Can we add a utils method for this?
	//Done
	if zkErr := utils.ValidateOffsetLimitValue(offset, limit); zkErr != nil {
		return response, zkErr
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
	//R: Can we moved to a utils method.
	//Done
	if zkErr := utils.ValidateOffsetLimitValue(offset, limit); zkErr != nil {
		return response, zkErr
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
			return response, &zkErr
		}

		return response, nil
	}

	zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorDbError, nil)
	return response, &zkErr
}
