package entity

import (
	"errors"
	"regexp"
	"unicode"
)

var (
	ErrPasswordTooShort  = errors.New("senha deve ter pelo menos 8 caracteres")
	ErrPasswordTooSimple = errors.New("senha deve conter pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial")
	ErrEmailInvalid      = errors.New("email inválido")
	ErrNameInvalid       = errors.New("nome deve conter apenas letras, números e espaços")
)

// validatePassword verifica se a senha atende aos requisitos mínimos de segurança
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
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

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPasswordTooSimple
	}

	return nil
}

// ValidateEmail verifica se o email é válido usando regex
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrEmailInvalid
	}
	return nil
}

// ValidateName verifica se o nome contém apenas caracteres permitidos
func ValidateName(name string) error {
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ú0-9\s]{2,100}$`)
	if !nameRegex.MatchString(name) {
		return ErrNameInvalid
	}
	return nil
}
