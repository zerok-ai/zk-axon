package errors

import (
	"github.com/kataras/iris/v12"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
)

var (
	ZkErrorBadRequestIssueHashEmpty       = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "IssueHash cannot be empty"}
	ZkErrorBadRequestScenarioIdEmpty      = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Scenario Id cannot be empty"}
	ZkErrorBadRequestTraceIdIdEmpty       = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "TraceId cannot be empty"}
	ZkErrorBadRequestErrorIdListIdEmpty   = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "ErrorId List cannot be empty"}
	ZkErrorBadRequestScenarioIdListEmpty  = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Scenario Id List cannot be empty"}
	ZkErrorBadRequestStartTimeEmpty       = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Start time cannot be empty"}
	ZkErrorBadRequestScenarioIdNotInteger = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Scenario Id is not integer"}
	ZkErrorBadRequestStartTimeNotNegative = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "Start time should be negative value"}
	ZkErrorBadRequestSpanIdEmpty          = zkErrors.ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "SpanId cannot be empty"}
	ZkErrorNotFound                       = zkErrors.ZkErrorType{Status: iris.StatusNotFound, Type: "NOT_FOUND", Message: "Not found"}
)
