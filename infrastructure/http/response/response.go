package response

import (
	"encoding/json"
	"net/http"
)

// ErrorCode representa um código de erro específico do domínio
type ErrorCode string

const (
	// Códigos de erro de autenticação
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired ErrorCode = "TOKEN_EXPIRED"

	// Códigos de erro de validação
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"

	// Códigos de erro de recursos
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeUnavailable   ErrorCode = "RESOURCE_UNAVAILABLE"

	// Códigos de erro de concorrência
	ErrCodeConcurrency ErrorCode = "CONCURRENT_MODIFICATION"
	ErrCodeConflict    ErrorCode = "CONFLICT"

	// Códigos de erro do sistema
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
)

// APIError representa um erro padronizado da API
type APIError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	HTTPStatus int         `json:"-"`
}

// APIResponse representa uma resposta padronizada da API
type APIResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    *APIError   `json:"error,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// NewErrorResponse cria uma nova resposta de erro
func NewErrorResponse(err *APIError) *APIResponse {
	return &APIResponse{
		Success: false,
		Error:   err,
	}
}

// NewSuccessResponse cria uma nova resposta de sucesso
func NewSuccessResponse(data interface{}, metadata interface{}) *APIResponse {
	return &APIResponse{
		Success:  true,
		Data:     data,
		Metadata: metadata,
	}
}

// JSON envia uma resposta JSON com o status HTTP apropriado
func (e *APIError) JSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.HTTPStatus)
	json.NewEncoder(w).Encode(NewErrorResponse(e))
}

// Error implementa a interface error
func (e *APIError) Error() string {
	return e.Message
}

// NewError cria um novo APIError com código de erro baseado no status HTTP
func NewError(status int, message string) *APIError {
	code := getErrorCodeFromStatus(status)
	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

// NewErrorWithDetails cria um novo APIError com detalhes adicionais
func NewErrorWithDetails(status int, message string, details interface{}) *APIError {
	err := NewError(status, message)
	err.Details = details
	return err
}

// getErrorCodeFromStatus mapeia códigos HTTP para códigos de erro do domínio
func getErrorCodeFromStatus(status int) ErrorCode {
	switch status {
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusBadRequest:
		return ErrCodeInvalidInput
	case http.StatusConflict:
		return ErrCodeConflict
	default:
		return ErrCodeInternal
	}
}
