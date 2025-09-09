package response

import (
	"net/http"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/repository"
)

// ErrorMapper mapeia erros do domínio para erros da API
func MapError(err error, requestID string) *APIError {
	if err == nil {
		return nil
	}

	switch err {
	// Erros de autenticação
	case entity.ErrUnauthorized:
		return &APIError{
			Code:       ErrCodeUnauthorized,
			Message:    "Não autorizado",
			RequestID:  requestID,
			HTTPStatus: http.StatusUnauthorized,
		}

	// Erros de usuário
	case entity.ErrUserNotFound:
		return &APIError{
			Code:       ErrCodeNotFound,
			Message:    "Usuário não encontrado",
			RequestID:  requestID,
			HTTPStatus: http.StatusNotFound,
		}
	case entity.ErrUserAlreadyExists:
		return &APIError{
			Code:       ErrCodeAlreadyExists,
			Message:    "Usuário já existe",
			RequestID:  requestID,
			HTTPStatus: http.StatusConflict,
		}
	case entity.ErrUserBlocked:
		return &APIError{
			Code:       ErrCodeForbidden,
			Message:    "Usuário bloqueado",
			RequestID:  requestID,
			HTTPStatus: http.StatusForbidden,
		}

	// Erros de chave
	case entity.ErrKeyNotFound:
		return &APIError{
			Code:       ErrCodeNotFound,
			Message:    "Chave não encontrada",
			RequestID:  requestID,
			HTTPStatus: http.StatusNotFound,
		}
	case entity.ErrKeyAlreadyExists:
		return &APIError{
			Code:       ErrCodeAlreadyExists,
			Message:    "Chave já existe",
			RequestID:  requestID,
			HTTPStatus: http.StatusConflict,
		}
	case entity.ErrKeyReserved:
		return &APIError{
			Code:       ErrCodeUnavailable,
			Message:    "Chave já está reservada",
			RequestID:  requestID,
			HTTPStatus: http.StatusConflict,
		}

	// Erros de reserva
	case entity.ErrReservationNotFound:
		return &APIError{
			Code:       ErrCodeNotFound,
			Message:    "Reserva não encontrada",
			RequestID:  requestID,
			HTTPStatus: http.StatusNotFound,
		}
	case entity.ErrReservationOverdue:
		return &APIError{
			Code:       ErrCodeConflict,
			Message:    "Reserva em atraso",
			RequestID:  requestID,
			HTTPStatus: http.StatusConflict,
		}

	// Erros de concorrência
	case repository.ErrConcurrentModification:
		return &APIError{
			Code:       ErrCodeConcurrency,
			Message:    "Modificação concorrente detectada",
			RequestID:  requestID,
			HTTPStatus: http.StatusConflict,
		}

	// Erro padrão para casos não mapeados
	default:
		return &APIError{
			Code:       ErrCodeInternal,
			Message:    "Erro interno do servidor",
			RequestID:  requestID,
			HTTPStatus: http.StatusInternalServerError,
		}
	}
}
