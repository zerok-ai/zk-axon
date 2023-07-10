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
	Protocol        string `json:"protocol"`
	RequestPayload  string `json:"request_payload"`
	ResponsePayload string `json:"response_payload"`
}

func ConvertSpanRawDataToSpanRawDataResponse(t []dto.SpanRawDataDetailsDto) (*SpanRawDataResponse, *error) {
	respMap := make(map[string]SpanRawDataDetails, 0)
	for _, v := range t {
		reqDecompressedStr, err := crypto.DecompressStringGzip(v.RequestPayload)
		if err != nil {
			return nil, &err
		}

		resDecompressedStr, err := crypto.DecompressStringGzip(v.ResponsePayload)
		if err != nil {
			zkLogger.Error(LogTag, err)
			return nil, &err
		}

		s := SpanRawDataDetails{
			Protocol:        v.Protocol,
			RequestPayload:  reqDecompressedStr,
			ResponsePayload: resDecompressedStr,
		}

		respMap[v.SpanId] = s
	}

	resp := SpanRawDataResponse{Spans: respMap}

	return &resp, nil
}
