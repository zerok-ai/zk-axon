package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Path Params
	IssueHash  = "issueHash"
	IncidentId = "incidentId"
	SpanId     = "spanId"
	ScenarioId = "scenarioId"
	//Scenario   = "scenarioId"

	// Query Params
	SpanIdQueryParam         = "span_id"
	ServicesQueryParam       = "services"
	ScenarioIdListQueryParam = "scenario_id_list"
	LimitQueryParam          = "limit"
	OffsetQueryParam         = "offset"
	StartTimeQueryParam      = "st"
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
