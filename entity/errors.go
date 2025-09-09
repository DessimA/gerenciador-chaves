package entity

import "errors"

var (
	// Erros de chave
	ErrKeyNotFound        = errors.New("chave não encontrada")
	ErrKeyAlreadyExists   = errors.New("chave já existe")
	ErrKeyInactive        = errors.New("chave está inativa")
	ErrKeyReserved        = errors.New("chave já está reservada")
	ErrKeyAlreadyReserved = errors.New("chave já está reservada por outro usuário")

	// Erros de usuário
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrUserAlreadyExists  = errors.New("usuário já existe")
	ErrUserBlocked        = errors.New("usuário está bloqueado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUnauthorized       = errors.New("não autorizado")

	// Erros de reserva
	ErrReservationNotFound      = errors.New("reserva não encontrada")
	ErrReservationOverdue       = errors.New("reserva em atraso")
	ErrReservationAlreadyExists = errors.New("usuário já possui reserva ativa")
	ErrInvalidReservationTime   = errors.New("tempo de reserva inválido")
	ErrCannotExtendReservation  = errors.New("não é possível estender a reserva")
)
