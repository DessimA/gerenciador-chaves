// Pacote models define as estruturas de dados da aplicação.
package models

import "time"

// Key representa a estrutura de uma chave no sistema.
type Key struct {
	ID              int        `json:"id"`
	ApartmentNumber string     `json:"apartment_number"`
	KeyType         string     `json:"key_type"`
	Status          string     `json:"status"` // "disponivel", "emprestada"
	BorrowedAt      *time.Time `json:"borrowed_at,omitempty"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty"`
	BorrowerName    *string    `json:"borrower_name,omitempty"`
}