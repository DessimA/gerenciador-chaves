package middleware

import (
	"github.com/dessima/gerenciador-chaves-api/infrastructure/http/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware para tratar erros de forma consistente
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Processa apenas se houver erros
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Obtém o ID da requisição do contexto
			requestID, exists := c.Get("requestID")
			requestIDStr := ""
			if exists {
				requestIDStr = requestID.(string)
			}

			// Mapeia o erro para um APIError
			apiError := response.MapError(err.Err, requestIDStr)

			// Se for um erro de validação, adiciona os detalhes
			if validationErrs, ok := err.Meta.(gin.H); ok && apiError.Code == response.ErrCodeValidation {
				apiError.Details = validationErrs
			}

			// Retorna a resposta de erro padronizada
			c.JSON(apiError.HTTPStatus, response.NewErrorResponse(apiError))
			c.Abort()
		}
	}
}
