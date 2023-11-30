package request

import "axon/internal/integrations/dto"

type UnsavedIntegrationRequestBody struct {
	Url string `json:"url"`
	dto.Authentication
}
