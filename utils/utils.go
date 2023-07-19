package utils

import (
	"time"
)

const (
	// Path Params
	IssueHash  = "issueHash"
	IncidentId = "incidentId"
	SpanId     = "spanId"
	//Scenario   = "scenarioId"

	// Query Params
	SpanIdQueryParam   = "span_id"
	ServicesQueryParam = "services"
	LimitQueryParam    = "limit"
	OffsetQueryParam   = "offset"
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
