package errors

import (
	"github.com/kataras/iris/v12"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

var (
	ZkErrorBadRequestIssueIdEmpty     = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "IssueHash cannot be empty"}
	ZkErrorBadRequestTraceIdIdEmpty   = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "TraceId cannot be empty"}
	ZkErrorBadRequestSpanIdEmpty      = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "SpanId cannot be empty"}
	ZkErrorBadRequestSourceEmpty      = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Source cannot be empty"}
	ZkErrorBadRequestDestinationEmpty = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Destination cannot be empty"}
)
