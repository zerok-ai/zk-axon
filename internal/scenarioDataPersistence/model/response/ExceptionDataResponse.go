package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"github.com/zerok-ai/zk-utils-go/crypto"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
)

type ExceptionDataResponse struct {
	Exceptions ExceptionDataDetailsMap `json:"exception_data_details"`
}

type ExceptionDataDetailsMap map[string]ExceptionDataDetails

type ExceptionDataDetails struct {
	Id            string `json:"id"`
	ExceptionBody string `json:"exception_body"`
}

func ConvertExceptionDataToExceptionDataResponse(t []dto.ExceptionTableDto) (ExceptionDataResponse, *error) {
	respMap := make(map[string]ExceptionDataDetails)
	resp := ExceptionDataResponse{}

	for _, v := range t {
		var exceptionDecompressedStr string
		var err error
		if v.ExceptionBody != nil && len(v.ExceptionBody) != 0 {
			exceptionDecompressedStr, err = crypto.DecompressStringGzip(v.ExceptionBody)
			if err != nil {
				zkLogger.Error(LogTag, "error decompressing exception body", err)
				return resp, &err
			}
		}

		s := ExceptionDataDetails{
			Id:            v.Id,
			ExceptionBody: exceptionDecompressedStr,
		}

		respMap[v.Id] = s
	}

	resp.Exceptions = respMap
	return resp, nil
}
