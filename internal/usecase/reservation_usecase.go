package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/portaria-keys/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReservationUseCase implementa os casos de uso relacionados às reservas
type ReservationUseCase struct {
	reservationRepo ReservationRepository
	keyRepo         KeyRepository
	userRepo        UserRepository
}

// NewReservationUseCase cria uma nova instância do ReservationUseCase
func NewReservationUseCase(
	reservationRepo ReservationRepository,
	keyRepo KeyRepository,
	userRepo UserRepository,
) *ReservationUseCase {
	return &ReservationUseCase{
		reservationRepo: reservationRepo,
		keyRepo:         keyRepo,
		userRepo:        userRepo,
	}
}

// TODO: Implementar todos os métodos do ReservationUseCase:

// CreateReservation cria uma nova reserva
func (uc *ReservationUseCase) CreateReservation(ctx context.Context, reservation *entity.Reservation) error {
	if err := reservation.ValidateForCreation(); err != nil {
		return err
	}

	user, err := uc.userRepo.GetByID(ctx, reservation.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}
	if user.IsBlocked {
		return entity.ErrUserBlocked
	}

	key, err := uc.keyRepo.GetByID(ctx, reservation.KeyID)
	if err != nil {
		return err
	}
	if key == nil {
		return entity.ErrKeyNotFound
	}
	if !key.CanBeReserved() {
		return entity.ErrKeyInactive
	}

	// Check if key is already reserved
	activeKeyReservation, err := uc.reservationRepo.GetActiveReservationByKey(ctx, reservation.KeyID)
	if err != nil && err != entity.ErrReservationNotFound {
		return err
	}
	if activeKeyReservation != nil {
		return entity.ErrKeyReserved
	}

	// Check if user already has an active reservation
	userReservations, err := uc.reservationRepo.GetByUserID(ctx, reservation.UserID)
	if err != nil {
		return err
	}
	for _, res := range userReservations {
		if res.Status == entity.ReservationStatusActive {
			return entity.ErrReservationAlreadyExists
		}
	}

	return uc.reservationRepo.Create(ctx, reservation)
}

// ReturnKey registra a devolução de uma chave
func (uc *ReservationUseCase) ReturnKey(ctx context.Context, reservationID primitive.ObjectID, userID primitive.ObjectID) error {
	reservation, err := uc.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}
	if reservation == nil {
		return entity.ErrReservationNotFound
	}

	// Check if the reservation belongs to the user or if the user is an admin
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	if reservation.UserID != userID && !user.IsAdmin() {
		return entity.ErrUnauthorized
	}

	if reservation.Status != entity.ReservationStatusActive {
		return errors.New("reservation is not active")
	}

	if err := reservation.MarkAsReturned(); err != nil {
		return err
	}

	// If overdue, block the user
	if reservation.IsOverdue() {
		if err := uc.userRepo.BlockUser(ctx, reservation.UserID); err != nil {
			return fmt.Errorf("failed to block user %s: %w", reservation.UserID.Hex(), err)
		}
	}

	return uc.reservationRepo.Update(ctx, reservation)
}

// GetUserReservations lista reservas de um usuário
func (uc *ReservationUseCase) GetUserReservations(ctx context.Context, userID primitive.ObjectID) ([]*entity.Reservation, error) {
	reservations, err := uc.reservationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Optionally, enrich reservations with Key and User details if needed for display
	// This would involve fetching Key and User entities for each reservation
	// For now, returning raw reservations as per the interface

	return reservations, nil
}

// GetAllReservations lista todas as reservas (apenas admin)
func (uc *ReservationUseCase) GetAllReservations(ctx context.Context, userRole entity.UserRole) ([]*entity.Reservation, error) {
	if userRole != entity.UserRoleAdmin {
		return nil, entity.ErrUnauthorized
	}

	reservations, err := uc.reservationRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Optionally, enrich reservations with Key and User details if needed for display
	// This would involve fetching Key and User entities for each reservation
	// For now, returning raw reservations as per the interface

	return reservations, nil
}

// ExtendReservation estende o prazo de uma reserva (apenas admin)
func (uc *ReservationUseCase) ExtendReservation(ctx context.Context, reservationID primitive.ObjectID, newDueTime time.Time, adminRole entity.UserRole) error {
	if adminRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	reservation, err := uc.reservationRepo.GetByID(ctx, reservationID)
	if err != nil {
		return err
	}
	if reservation == nil {
		return entity.ErrReservationNotFound
	}

	if !reservation.CanBeExtended() {
		return entity.ErrCannotExtendReservation
	}

	if newDueTime.Before(reservation.DueAt) {
		return errors.New("new due time cannot be before current due time")
	}

	reservation.DueAt = newDueTime
	reservation.UpdatedAt = time.Now()

	return uc.reservationRepo.Update(ctx, reservation)
}

// ProcessOverdueReservations processa reservas em atraso (job automático)
func (uc *ReservationUseCase) ProcessOverdueReservations(ctx context.Context) error {
	overdueReservations, err := uc.reservationRepo.GetOverdueReservations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get overdue reservations: %w", err)
	}

	for _, res := range overdueReservations {
		if res.Status == entity.ReservationStatusActive && time.Now().After(res.DueAt) {
			// Mark as overdue
			if err := res.MarkAsOverdue(); err != nil {
				return fmt.Errorf("failed to mark reservation %s as overdue: %w", res.ID.Hex(), err)
			}
			if err := uc.reservationRepo.Update(ctx, res); err != nil {
				return fmt.Errorf("failed to update reservation %s: %w", res.ID.Hex(), err)
			}

			// Block user if overdue
			if err := uc.userRepo.BlockUser(ctx, res.UserID); err != nil {
				return fmt.Errorf("failed to block user %s: %w", res.UserID.Hex(), err)
			}
		}
	}

	return nil
}