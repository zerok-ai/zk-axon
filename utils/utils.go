package utils

import (
	"fmt"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
	"strings"
	"time"
)

const (
	// Path Params
	IssueHash                    = "issueHash"
	IncidentId                   = "incidentId"
	SpanId                       = "spanId"
	ScenarioId                   = "scenarioId"
	IntegrationIdxPathParam      = "integrationId"
	MetricAttributeNamePathParam = "metric"

	Namespace     = "namespace"
	TraceId       = "traceId"
	PodId         = "podId"
	PromQuery     = "promQuery"
	IntegrationId = "integrationId"

	RateIntervalQueryParam = "rate-interval"
	TimeQueryParam         = "time"

	// Query Params
	IssueHashQueryParam      = "issue_hash"
	SpanIdQueryParam         = "span_id"
	ServicesQueryParam       = "services"
	ScenarioIdListQueryParam = "scenario_id_list"
	LimitQueryParam          = "limit"
	OffsetQueryParam         = "offset"
	StartTimeQueryParam      = "st"
	EndTimeQueryParam        = "et"
	DurationQueryParam       = "duration"

	LogTag = "utils"

	StatusError   = "error"
	StatusSuccess = "success"

	ConnectionSuccessful = "Connection Successful"
	ConnectionFailed     = "Connection Failed"
)

func CalendarDaysBetween(start, end time.Time) int {
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	return days
}

func HoursBetween(start, end time.Time) int {
	duration := end.Sub(start)
	hours := int(duration.Hours())
	return hours
}

func ParseTimeString(input string) (time.Duration, error) {
	var duration time.Duration
	var multiplier time.Duration

	switch {
	case strings.HasSuffix(input, "m"):
		multiplier = time.Minute
	case strings.HasSuffix(input, "h"):
		multiplier = time.Hour
	case strings.HasSuffix(input, "d"):
		multiplier = 24 * time.Hour
	default:
		return 0, fmt.Errorf("unsupported input format")
	}

	numericPart := strings.TrimSuffix(input, string(input[len(input)-1]))
	val, err := strconv.Atoi(numericPart)
	if err != nil {
		return 0, err
	}

	duration = time.Duration(val) * multiplier
	return duration, nil
}

func ValidateOffsetLimitValue(offset, limit int) *zkErrors.ZkError {
	//R: Can we also add max number for limit, if the dashboard or a bad actor sends a very high limit,
	//then it might end up crashing the service.
	if offset < 0 || limit < 1 {
		zkErr := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, nil)
		zkLogger.Error(LogTag, fmt.Sprintf("value of limit or offset is invalid, limit: %d, offset: %d", limit, offset), zkErr)
		return &zkErr
	}

	return nil
}

func MergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func BuildZkError(logTag string, errStr ...string) *zkErrors.ZkError {
	err := zkErrors.ZkErrorBuilder{}.Build(zkErrors.ZkErrorBadRequest, strings.Join(errStr, " "))
	zkLogger.Error(logTag, errStr)
	return &err
}
