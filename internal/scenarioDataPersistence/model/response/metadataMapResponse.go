package scenariodataresponse

import (
	"axon/internal/scenarioDataPersistence/model/dto"
	"strings"
)

type MetadataMapResponse struct {
	MetadataMapList []dto.MetadataMapDto `json:"metadata_map_list"`
}

func ConvertMetadataMapToMetadataMapResponse(t []dto.MetadataMapDto) (*MetadataMapResponse, *error) {
	var resList []dto.MetadataMapDto
	for _, v := range t {
		x := dto.MetadataMapDto{
			Source:       removeLastTwoStrings(v.Source),
			Destination:  removeLastTwoStrings(v.Destination),
			TraceCount:   v.TraceCount,
			ProtocolList: v.ProtocolList,
		}
		resList = append(resList, x)
	}
	resp := MetadataMapResponse{MetadataMapList: resList}
	return &resp, nil
}

func removeLastTwoStrings(input string) string {
	parts := strings.Split(input, "-")
	if len(parts) <= 2 {
		return input
	}

	trimmedParts := parts[:len(parts)-2]
	return strings.Join(trimmedParts, "-")
}
