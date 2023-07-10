package utils

import (
	"time"
)

const (
	// Path Params
	IssueId    = "issueId"
	IncidentId = "incidentId"
	SpanId     = "spanId"

	// Query Params
	SpanIdQueryParam      = "span_id"
	SourceQueryParam      = "source"
	DestinationQueryParam = "destination"
	LimitQueryParam       = "limit"
	OffsetQueryParam      = "offset"
)

func CalendarDaysBetween(start, end time.Time) int {
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	return days
}
