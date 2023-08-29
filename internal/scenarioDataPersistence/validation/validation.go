package validation

import (
	zkErrorsAxon "axon/utils/zkerrors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

func GetIssuesListWithDetails(offset, limit, startTime string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(startTime) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestStartTimeEmpty, nil)
		return &zkErr
	}

	if zkErr := ValidateLimit(limit); zkErr != nil {
		return zkErr
	}

	if zkErr := ValidateOffset(offset); zkErr != nil {
		return zkErr
	}

	return nil
}

func ValidateGetScenarioDetails(scenarioIds, startTime string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(startTime) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestStartTimeEmpty, nil)
		return &zkErr
	}

	if zkCommon.IsEmpty(scenarioIds) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestScenarioIdListEmpty, nil)
		return &zkErr
	}

	return nil
}

func ValidateIssueHashOffsetAndLimit(issueHash, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(issueHash) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestIssueHashEmpty, nil)
		return &zkErr
	}

	if zkErr := ValidateLimit(limit); zkErr != nil {
		return zkErr
	}

	if zkErr := ValidateOffset(offset); zkErr != nil {
		return zkErr
	}

	return nil
}

func ValidateScenarioIdOffsetAndLimit(scenarioId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(scenarioId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestScenarioIdEmpty, nil)
		return &zkErr
	}

	if zkErr := ValidateLimit(limit); zkErr != nil {
		return zkErr
	}

	if zkErr := ValidateOffset(offset); zkErr != nil {
		return zkErr
	}

	return nil
}

func ValidateGetSpanRawDataApi(traceId, spanId string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(traceId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestTraceIdIdEmpty, nil)
		return &zkErr
	}

	if zkCommon.IsEmpty(spanId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestSpanIdEmpty, nil)
		return &zkErr
	}

	return nil
}

func ValidateGetIncidentDetailsApi(traceId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(traceId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsAxon.ZkErrorBadRequestTraceIdIdEmpty, nil)
		return &zkErr
	}

	if zkErr := ValidateLimit(limit); zkErr != nil {
		return zkErr
	}

	if zkErr := ValidateOffset(offset); zkErr != nil {
		return zkErr
	}

	return nil
}

func ValidateLimit(limit string) *zkerrors.ZkError {
	if !zkCommon.IsEmpty(limit) {
		//R: I think we can also check if limit is greater than 0.
		//R: Can we also add max number for limit, if the dashboard or a bad actor sends a very high limit,
		//then it might end up crashing the service.
		_, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}
	return nil
}

func ValidateOffset(offset string) *zkerrors.ZkError {
	if !zkCommon.IsEmpty(offset) {
		//R: I think we can also check if offset is greater than 0.
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}
	return nil
}
