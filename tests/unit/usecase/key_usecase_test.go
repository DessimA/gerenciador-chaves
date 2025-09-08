package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/portaria-keys/internal/entity"
	"github.com/portaria-keys/internal/usecase"
	"github.com/portaria-keys/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestKeyUseCase_CreateKey(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository) // Not used in CreateKey, but needed for NewKeyUseCase
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()

	// Test case 1: Admin user, successful creation
	key := entity.NewKey("New Key", "Description")
	mockKeyRepo.On("GetByName", ctx, key.Name).Return(nil, entity.ErrKeyNotFound).Once()
	mockKeyRepo.On("Create", ctx, mock.AnythingOfType("*entity.Key")).Return(nil).Once()

	err := keyUseCase.CreateKey(ctx, key, entity.UserRoleAdmin)
	assert.Nil(t, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	key = entity.NewKey("Another Key", "Description")
	err = keyUseCase.CreateKey(ctx, key, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: Invalid key data
	key = entity.NewKey("", "Description") // Invalid name
	err = keyUseCase.CreateKey(ctx, key, entity.UserRoleAdmin)
	assert.NotNil(t, err)

	// Test case 4: Key with same name already exists
	key = entity.NewKey("Existing Key", "Description")
	mockKeyRepo.On("GetByName", ctx, key.Name).Return(key, nil).Once()

	err = keyUseCase.CreateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrKeyAlreadyExists, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 5: Repository error during GetByName
	key = entity.NewKey("Repo Error Key", "Description")
	repoError := errors.New("database error")
	mockKeyRepo.On("GetByName", ctx, key.Name).Return(nil, repoError).Once()

	err = keyUseCase.CreateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 6: Repository error during Create
	key = entity.NewKey("Create Error Key", "Description")
	mockKeyRepo.On("GetByName", ctx, key.Name).Return(nil, entity.ErrKeyNotFound).Once()
	mockKeyRepo.On("Create", ctx, mock.AnythingOfType("*entity.Key")).Return(repoError).Once()

	err = keyUseCase.CreateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockKeyRepo.AssertExpectations(t)
}

func TestKeyUseCase_GetKeyByID(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()

	// Test case 1: Key found
	expectedKey := &entity.Key{ID: id, Name: "Test Key"}
	mockKeyRepo.On("GetByID", ctx, id).Return(expectedKey, nil).Once()

	key, err := keyUseCase.GetKeyByID(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, expectedKey, key)
	mockKeyRepo.AssertExpectations(t)

	// Test case 2: Key not found
	mockKeyRepo.On("GetByID", ctx, id).Return(nil, entity.ErrKeyNotFound).Once()

	key, err = keyUseCase.GetKeyByID(ctx, id)
	assert.Equal(t, entity.ErrKeyNotFound, err)
	assert.Nil(t, key)
	mockKeyRepo.AssertExpectations(t)

	// Test case 3: Repository error
	repoError := errors.New("database error")
	mockKeyRepo.On("GetByID", ctx, id).Return(nil, repoError).Once()

	key, err = keyUseCase.GetKeyByID(ctx, id)
	assert.Equal(t, repoError, err)
	assert.Nil(t, key)
	mockKeyRepo.AssertExpectations(t)
}

func TestKeyUseCase_GetAllKeys(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()

	// Test case 1: Successful retrieval
	expectedKeys := []*entity.Key{{Name: "Key1"}, {Name: "Key2"}}
	mockKeyRepo.On("GetAll", ctx).Return(expectedKeys, nil).Once()

	keys, err := keyUseCase.GetAllKeys(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedKeys, keys)
	mockKeyRepo.AssertExpectations(t)

	// Test case 2: Repository error
	repoError := errors.New("database error")
	mockKeyRepo.On("GetAll", ctx).Return(nil, repoError).Once()

	keys, err = keyUseCase.GetAllKeys(ctx)
	assert.Equal(t, repoError, err)
	assert.Nil(t, keys)
	mockKeyRepo.AssertExpectations(t)
}

func TestKeyUseCase_UpdateKey(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()

	// Test case 1: Admin user, successful update
	existingKey := &entity.Key{ID: id, Name: "Old Name", Description: "Old Desc", IsActive: true}
	updatedKey := &entity.Key{ID: id, Name: "New Name", Description: "New Desc", IsActive: false}
	mockKeyRepo.On("GetByID", ctx, id).Return(existingKey, nil).Once()
	mockKeyRepo.On("Update", ctx, mock.AnythingOfType("*entity.Key")).Return(nil).Once()

	err := keyUseCase.UpdateKey(ctx, updatedKey, entity.UserRoleAdmin)
	assert.Nil(t, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	key := &entity.Key{ID: id, Name: "New Name"}
	err = keyUseCase.UpdateKey(ctx, key, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: Key not found
	key = &entity.Key{ID: id, Name: "New Name"}
	mockKeyRepo.On("GetByID", ctx, id).Return(nil, entity.ErrKeyNotFound).Once()

	err = keyUseCase.UpdateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrKeyNotFound, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 4: Repository error during GetByID
	key = &entity.Key{ID: id, Name: "New Name"}
	repoError := errors.New("database error")
	mockKeyRepo.On("GetByID", ctx, id).Return(nil, repoError).Once()

	err = keyUseCase.UpdateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 5: Repository error during Update
	key = &entity.Key{ID: id, Name: "New Name"}
	mockKeyRepo.On("GetByID", ctx, id).Return(existingKey, nil).Once()
	mockKeyRepo.On("Update", ctx, mock.AnythingOfType("*entity.Key")).Return(repoError).Once()

	err = keyUseCase.UpdateKey(ctx, key, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockKeyRepo.AssertExpectations(t)
}

func TestKeyUseCase_DeleteKey(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()
	id := primitive.NewObjectID()

	// Test case 1: Admin user, successful deletion
	key := &entity.Key{ID: id, Name: "Test Key"}
	mockKeyRepo.On("GetByID", ctx, id).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, id).Return(nil, entity.ErrReservationNotFound).Once()
	mockKeyRepo.On("Delete", ctx, id).Return(nil).Once()

	err := keyUseCase.DeleteKey(ctx, id, entity.UserRoleAdmin)
	assert.Nil(t, err)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	err = keyUseCase.DeleteKey(ctx, id, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: Key not found
	mockKeyRepo.On("GetByID", ctx, id).Return(nil, entity.ErrKeyNotFound).Once()

	err = keyUseCase.DeleteKey(ctx, id, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrKeyNotFound, err)
	mockKeyRepo.AssertExpectations(t)

	// Test case 4: Key has active reservations
	activeReservation := &entity.Reservation{KeyID: id, Status: entity.ReservationStatusActive}
	mockKeyRepo.On("GetByID", ctx, id).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, id).Return(activeReservation, nil).Once()

	err = keyUseCase.DeleteKey(ctx, id, entity.UserRoleAdmin)
	assert.EqualError(t, err, "cannot delete key with active reservations")
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 5: Repository error during Delete
	repoError := errors.New("database error")
	mockKeyRepo.On("GetByID", ctx, id).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, id).Return(nil, entity.ErrReservationNotFound).Once()
	mockKeyRepo.On("Delete", ctx, id).Return(repoError).Once()

	err = keyUseCase.DeleteKey(ctx, id, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)
}

func TestKeyUseCase_GetAvailableKeys(t *testing.T) {
	mockKeyRepo := new(mocks.KeyRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	keyUseCase := usecase.NewKeyUseCase(mockKeyRepo, mockReservationRepo)

	ctx := context.Background()

	// Test case 1: Successful retrieval
	expectedKeys := []*entity.Key{{Name: "Available1"}, {Name: "Available2"}}
	mockKeyRepo.On("GetAvailableKeys", ctx).Return(expectedKeys, nil).Once()

	keys, err := keyUseCase.GetAvailableKeys(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectedKeys, keys)
	mockKeyRepo.AssertExpectations(t)

	// Test case 2: Repository error
	repoError := errors.New("database error")
	mockKeyRepo.On("GetAvailableKeys", ctx).Return(nil, repoError).Once()

	keys, err = keyUseCase.GetAvailableKeys(ctx)
	assert.Equal(t, repoError, err)
	assert.Nil(t, keys)
	mockKeyRepo.AssertExpectations(t)
}
