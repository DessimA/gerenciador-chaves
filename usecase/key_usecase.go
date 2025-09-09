package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyUseCase implementa os casos de uso relacionados às chaves
type KeyUseCase struct {
	keyRepo KeyRepository
	reservationRepo ReservationRepository
}

// NewKeyUseCase cria uma nova instância do KeyUseCase
func NewKeyUseCase(keyRepo KeyRepository, reservationRepo ReservationRepository) *KeyUseCase {
	return &KeyUseCase{
		keyRepo: keyRepo,
		reservationRepo: reservationRepo,
	}
}

// TODO: Implementar todos os métodos do KeyUseCase:

// CreateKey cria uma nova chave (apenas admin)
func (uc *KeyUseCase) CreateKey(ctx context.Context, key *entity.Key, userRole entity.UserRole) error {
	if userRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	if err := key.ValidateForCreation(); err != nil {
		return err
	}

	existingKey, err := uc.keyRepo.GetByName(ctx, key.Name)
	if err != nil && err != entity.ErrKeyNotFound {
		return err
	}
	if existingKey != nil {
		return entity.ErrKeyAlreadyExists
	}

	key.CreatedAt = time.Now()
	key.UpdatedAt = time.Now()

	return uc.keyRepo.Create(ctx, key)
}

// GetKeyByID busca uma chave por ID
func (uc *KeyUseCase) GetKeyByID(ctx context.Context, id primitive.ObjectID) (*entity.Key, error) {
	key, err := uc.keyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, entity.ErrKeyNotFound
	}
	return key, nil
}

// GetAllKeys lista todas as chaves
func (uc *KeyUseCase) GetAllKeys(ctx context.Context) ([]*entity.Key, error) {
	keys, err := uc.keyRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// UpdateKey atualiza uma chave existente (apenas admin)
func (uc *KeyUseCase) UpdateKey(ctx context.Context, key *entity.Key, userRole entity.UserRole) error {
	if userRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	existingKey, err := uc.keyRepo.GetByID(ctx, key.ID)
	if err != nil {
		return err
	}
	if existingKey == nil {
		return entity.ErrKeyNotFound
	}

	// Update only allowed fields
	existingKey.Name = key.Name
	existingKey.Description = key.Description
	existingKey.IsActive = key.IsActive
	existingKey.UpdatedAt = time.Now()

	if err := existingKey.ValidateForUpdate(); err != nil {
		return err
	}

	return uc.keyRepo.Update(ctx, existingKey)
}

// DeleteKey remove uma chave (apenas admin)
func (uc *KeyUseCase) DeleteKey(ctx context.Context, id primitive.ObjectID, userRole entity.UserRole) error {
	if userRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	existingKey, err := uc.keyRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingKey == nil {
		return entity.ErrKeyNotFound
	}

	// Check for active reservations
	activeReservation, err := uc.reservationRepo.GetActiveReservationByKey(ctx, id)
	if err != nil && err != entity.ErrReservationNotFound {
		return err
	}
	if activeReservation != nil {
		return errors.New("cannot delete key with active reservations")
	}

	return uc.keyRepo.Delete(ctx, id)
}

// GetAvailableKeys lista chaves disponíveis para reserva
func (uc *KeyUseCase) GetAvailableKeys(ctx context.Context) ([]*entity.Key, error) {
	keys, err := uc.keyRepo.GetAvailableKeys(ctx)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
