package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"github.com/zerok-ai/zk-utils-go/crypto"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
)

type SpanRawDataResponse struct {
	Spans SpansRawDataDetailsMap `json:"span_raw_data_details"`
}

type SpansRawDataDetailsMap map[string]SpanRawDataDetails

type SpanRawDataDetails struct {
	TraceID     string `json:"trace_id"`
	SpanID      string `json:"span_id"`
	Protocol    string `json:"protocol"`
	ReqHeaders  string `json:"req_headers"`
	RespHeaders string `json:"resp_headers"`
	IsTruncated bool   `json:"is_truncated"`
	ReqBody     string `json:"req_body"`
	RespBody    string `json:"resp_body"`
}

func ConvertSpanRawDataToSpanRawDataResponse(t []dto.SpanRawDataDetailsDto) (SpanRawDataResponse, *error) {
	respMap := make(map[string]SpanRawDataDetails, 0)
	resp := SpanRawDataResponse{}

	for _, v := range t {
		var reqDecompressedStr, resDecompressedStr string
		var err error
		if v.ReqBody != nil && len(v.ReqBody) != 0 {
			reqDecompressedStr, err = crypto.DecompressStringGzip(v.ReqBody)
			if err != nil {
				zkLogger.Error(LogTag, "error decompressing request payload", err)
				return resp, &err
			}
		}

		if v.RespBody != nil && len(v.RespBody) != 0 {
			resDecompressedStr, err = crypto.DecompressStringGzip(v.RespBody)
			if err != nil {
				zkLogger.Error(LogTag, "error decompressing response payload", err)
				return resp, &err
			}
		}

		s := SpanRawDataDetails{
			TraceID:     v.TraceID,
			SpanID:      v.SpanID,
			Protocol:    v.Protocol,
			ReqHeaders:  v.ReqHeaders,
			RespHeaders: v.RespHeaders,
			IsTruncated: v.IsTruncated,
			ReqBody:     reqDecompressedStr,
			RespBody:    resDecompressedStr,
		}

		respMap[v.SpanID] = s
	}

	resp.Spans = respMap
	return resp, nil
}
