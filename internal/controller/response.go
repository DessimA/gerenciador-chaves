package controller

// APIError representa uma estrutura de erro padronizada para a API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}