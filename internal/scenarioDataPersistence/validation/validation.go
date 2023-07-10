package validation

import (
	zkErrorsScenarioManager "axon/utils/zkerrors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

func GetIssuesListWithDetails(source, destination, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(destination) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestDestinationEmpty, nil)
		return &zkErr
	}

	if zkCommon.IsEmpty(source) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestSourceEmpty, nil)
		return &zkErr
	}

	if !zkCommon.IsEmpty(limit) {
		_, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(offset) {
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestOffsetIsNotInteger, nil)
			return &zkErr
		}
	}

	return nil
}

func ValidateGetIncidentApi(issueId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(issueId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestIssueIdEmpty, nil)
		return &zkErr
	}

	if !zkCommon.IsEmpty(limit) {
		_, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(offset) {
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestOffsetIsNotInteger, nil)
			return &zkErr
		}
	}

	return nil
}

func ValidateGetSpanRawDataApi(traceId, spanId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(traceId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestTraceIdIdEmpty, nil)
		return &zkErr
	}

	if zkCommon.IsEmpty(spanId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestSpanIdEmpty, nil)
		return &zkErr
	}

	if !zkCommon.IsEmpty(limit) {
		_, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(offset) {
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestOffsetIsNotInteger, nil)
			return &zkErr
		}
	}

	return nil
}

func ValidateGetIncidentDetailsApi(traceId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(traceId) {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkErrorsScenarioManager.ZkErrorBadRequestTraceIdIdEmpty, nil)
		return &zkErr
	}

	if !zkCommon.IsEmpty(limit) {
		_, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}

	if !zkCommon.IsEmpty(offset) {
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestOffsetIsNotInteger, nil)
			return &zkErr
		}
	}

	return nil
}
