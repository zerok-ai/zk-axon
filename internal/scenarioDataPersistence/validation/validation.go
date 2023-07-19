package validation

import (
	zkErrorsAxon "axon/utils/zkerrors"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"strconv"
)

func GetIssuesListWithDetails(offset, limit string) *zkerrors.ZkError {
	if zkErr := ValidateLimit(limit); zkErr != nil {
		return zkErr
	}

	if zkErr := ValidateOffset(offset); zkErr != nil {
		return zkErr
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
		_, err := strconv.Atoi(offset)
		if err != nil {
			zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorBadRequestLimitIsNotInteger, nil)
			return &zkErr
		}
	}
	return nil
}
