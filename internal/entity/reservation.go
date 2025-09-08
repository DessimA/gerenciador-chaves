package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReservationStatus define os status de uma reserva
type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "active"
	ReservationStatusReturned ReservationStatus = "returned"
	ReservationStatusOverdue  ReservationStatus = "overdue"
)

// Reservation representa uma reserva de chave
type Reservation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	KeyID      primitive.ObjectID `bson:"key_id" json:"key_id" validate:"required"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	ReservedAt time.Time          `bson:"reserved_at" json:"reserved_at"`
	DueAt      time.Time          `bson:"due_at" json:"due_at" validate:"required"`
	ReturnedAt *time.Time         `bson:"returned_at,omitempty" json:"returned_at,omitempty"`
	Status     ReservationStatus  `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewReservation cria uma nova instância de Reservation com timestamps e status preenchidos
func NewReservation(keyID, userID primitive.ObjectID, dueAt time.Time) *Reservation {
	now := time.Now()
	return &Reservation{
		ID:         primitive.NewObjectID(),
		KeyID:      keyID,
		UserID:     userID,
		ReservedAt: now,
		DueAt:      dueAt,
		Status:     ReservationStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// ValidateForCreation valida os campos da Reservation para criação
func (r *Reservation) ValidateForCreation() error {
	if r.DueAt.Before(time.Now()) {
		return errors.New("due_at must be in the future")
	}
	return validate.Struct(r)
}

// ValidateReturnTime valida se o ReturnedAt é válido (não antes de ReservedAt)
func (r *Reservation) ValidateReturnTime() error {
	if r.ReturnedAt != nil && r.ReturnedAt.Before(r.ReservedAt) {
		return errors.New("returned_at cannot be before reserved_at")
	}
	return nil
}

// IsOverdue verifica se a reserva está em atraso
func (r *Reservation) IsOverdue() bool {
	return r.Status == ReservationStatusActive && time.Now().After(r.DueAt)
}

// CalculateOverdueTime calcula o tempo de atraso da reserva
func (r *Reservation) CalculateOverdueTime() time.Duration {
	if r.IsOverdue() {
		return time.Since(r.DueAt)
	}
	return 0
}

// CanBeExtended verifica se a reserva pode ter seu prazo estendido
func (r *Reservation) CanBeExtended() bool {
	return r.Status == ReservationStatusActive
}

// MarkAsReturned marca a reserva como devolvida
func (r *Reservation) MarkAsReturned() error {
	if r.Status != ReservationStatusActive {
		return errors.New("reservation is not active")
	}
	now := time.Now()
	r.ReturnedAt = &now
	r.Status = ReservationStatusReturned
	r.UpdatedAt = now
	return nil
}

// MarkAsOverdue marca a reserva como em atraso
func (r *Reservation) MarkAsOverdue() error {
	if r.Status != ReservationStatusActive {
		return errors.New("reservation is not active")
	}
	r.Status = ReservationStatusOverdue
	r.UpdatedAt = time.Now()
	return nil
}