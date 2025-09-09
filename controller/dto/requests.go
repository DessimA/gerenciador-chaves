package dto

// KeyRequest representa os dados necessários para criar ou atualizar uma chave
type KeyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsActive    bool   `json:"isActive"`
}

// ReservationRequest representa os dados necessários para criar uma reserva
type ReservationRequest struct {
	KeyID       string `json:"keyId" binding:"required"`
	ReturnDate  string `json:"returnDate" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// ExtendReservationRequest representa os dados necessários para estender uma reserva
type ExtendReservationRequest struct {
	NewReturnDate string `json:"newReturnDate" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
}
