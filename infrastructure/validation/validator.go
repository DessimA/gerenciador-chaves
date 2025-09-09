package validation

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Registra validações customizadas
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("name", validateName)
}

// ValidateRequest valida uma struct e retorna erros formatados para a API
func ValidateRequest(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		return formatBindError(err)
	}

	if err := validate.Struct(request); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// formatBindError formata erros de binding JSON
func formatBindError(err error) error {
	return &ValidationError{
		Field:   "request",
		Message: "JSON inválido",
		Details: err.Error(),
	}
}

// formatValidationError formata erros de validação
func formatValidationError(err error) error {
	var errors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(err.Field()),
			Message: getErrorMessage(err),
			Tag:     err.Tag(),
		})
	}

	return errors
}

// ValidationError representa um erro de validação individual
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Details string `json:"details,omitempty"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// ValidationErrors representa uma coleção de erros de validação
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	return e[0].Error()
}

// getErrorMessage retorna uma mensagem amigável para cada tipo de erro de validação
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Campo obrigatório"
	case "email":
		return "Email inválido"
	case "min":
		return "Valor muito curto"
	case "max":
		return "Valor muito longo"
	case "password":
		return "Senha não atende aos requisitos mínimos"
	case "name":
		return "Nome contém caracteres inválidos"
	default:
		return "Campo inválido"
	}
}
