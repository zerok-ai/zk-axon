package response

type GenericQueryResponse struct {
	Result interface{} `json:"result"`
	Type   string      `json:"type"`
}
