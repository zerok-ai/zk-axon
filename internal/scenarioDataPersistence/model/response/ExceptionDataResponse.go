package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"github.com/zerok-ai/zk-utils-go/crypto"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
)

type ErrorDataResponse struct {
	Errors []ErrorDataDetails `json:"errors"`
}

type ErrorDataDetailsMap map[string]ErrorDataDetails

type ErrorDataDetails struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func ConvertErrorDataToErrorDataResponse(t []dto.ErrorDataTableDto) (ErrorDataResponse, *error) {
	respList := make([]ErrorDataDetails, 0)
	var resp ErrorDataResponse

	for _, v := range t {
		var errorDecompressedStr string
		var err error
		if v.Data != nil && len(v.Data) != 0 {
			errorDecompressedStr, err = crypto.DecompressStringGzip(v.Data)
			if err != nil {
				zkLogger.Error(LogTag, "error decompressing error body", err)
				return resp, &err
			}
		}

		s := ErrorDataDetails{
			Id:   v.Id,
			Data: errorDecompressedStr,
		}

		respList = append(respList, s)
	}

	resp.Errors = respList
	return resp, nil
}
