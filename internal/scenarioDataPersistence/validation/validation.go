package validation

import (
	zkErrorsAxon "axon/utils/zkerrors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

var VALIDATE_LOG_TAG = "validation"

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

func ValidateIdStringOffsetAndLimit(scenarioId, offset, limit string) *zkerrors.ZkError {
	if zkCommon.IsEmpty(scenarioId) {
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
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
		if limitInt < 1 {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, "Limit is invalid.")
			logger.Debug(VALIDATE_LOG_TAG, "Limit is invalid.")
			return &zkErr
		}
	}
	return nil
}

func ValidateOffset(offset string) *zkerrors.ZkError {
	if !zkCommon.IsEmpty(offset) {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}

		if offsetInt < 1 {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequest, "Offset is invalid.")
			logger.Debug(VALIDATE_LOG_TAG, "Offset is invalid.")
			return &zkErr
		}
	}
	return nil
}
