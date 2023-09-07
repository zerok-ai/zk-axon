package request

import (
	"time"
)

type PodInfoRequest struct {
	Namespace    string
	Pod          string
	RateInterval string
	StartTime    time.Time
	EndTime      time.Time
	Timestamp    int64
}
