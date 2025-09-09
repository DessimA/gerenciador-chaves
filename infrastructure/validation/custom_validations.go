package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// validatePassword verifica se a senha atende aos requisitos mínimos
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Mínimo de 8 caracteres
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateName verifica se o nome contém apenas caracteres permitidos
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	// Nome deve ter entre 2 e 100 caracteres
	if len(name) < 2 || len(name) > 100 {
		return false
	}

	// Permite letras, números e espaços, com suporte a caracteres acentuados
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ú0-9\s]{2,100}$`)
	return nameRegex.MatchString(name)
}
