package contract

// ResponseJSON objeto de resposta de sucesso
type ResponseJSON struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
