package request

import (
	"time"
)

type PromRequestMeta struct {
	Namespace        string
	Pod              string
	RateInterval     string
	StartTime        time.Time
	EndTime          time.Time
	Timestamp        int64
	PodsListStr      string
	NamespaceListStr string
}
