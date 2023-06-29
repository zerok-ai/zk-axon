package utils

import (
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"regexp"
	"strings"
	"time"
)

const (
	ScenarioType = "scenario_type"
	ScenarioId   = "scenario_id"
	TraceId      = "trace_id"
	SpanId       = "span_id"
	Source       = "source"
	Destination  = "destination"
	Limit        = "limit"
	Offset       = "offset"
	Duration     = "duration"
)

var TimeUnitPxl = []string{"s", "m", "h", "d", "mon"}

func ParseTimestamp(timestamp string) (time.Time, error) {
	// Define the layout of the timestamp string
	layout := "2006-01-02T15:04:05.999999Z"

	// Parse the timestamp string using the specified layout
	parsedTime, err := time.Parse(layout, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func CalendarDaysBetween(start, end time.Time) int {
	start = start.Truncate(24 * time.Hour)
	end = end.Truncate(24 * time.Hour)
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	return days
}

func IsValidPxlTime(s string) bool {
	re := regexp.MustCompile("[0-9]+")
	d := re.FindAllString(s, -1)
	if len(d) != 1 {
		return false
	}

	t := strings.Split(s, d[0])
	var params = make([]string, 0)
	for _, v := range t {
		if !zkCommon.IsEmpty(v) {
			params = append(params, v)
		}
	}
	if len(params) == 2 {
		if !zkCommon.Contains(TimeUnitPxl, params[1]) || params[0] != "-" {
			return false
		}
	} else if len(params) == 1 {
		if !zkCommon.Contains(TimeUnitPxl, params[0]) {
			return false
		}
	} else {
		return false
	}

	return true
}
