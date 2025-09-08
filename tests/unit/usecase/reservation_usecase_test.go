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

func TestReservationUseCase_CreateReservation(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()
	keyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	dueAt := time.Now().Add(time.Hour)

	// Test case 1: Successful creation
	reservation := entity.NewReservation(keyID, userID, dueAt)
	user := &entity.User{ID: userID, IsBlocked: false}
	key := &entity.Key{ID: keyID, IsActive: true}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, keyID).Return(nil, entity.ErrReservationNotFound).Once()
	mockReservationRepo.On("GetByUserID", ctx, userID).Return([]*entity.Reservation{}, nil).Once()
	mockReservationRepo.On("Create", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err := reservationUseCase.CreateReservation(ctx, reservation)
	assert.Nil(t, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Invalid reservation data (DueAt in past)
	reservation = entity.NewReservation(keyID, userID, time.Now().Add(-time.Hour))
	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.NotNil(t, err)

	// Test case 3: User not found
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, entity.ErrUserNotFound).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 4: User is blocked
	userBlocked := &entity.User{ID: userID, IsBlocked: true}
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(userBlocked, nil).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrUserBlocked, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Key not found
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(nil, entity.ErrKeyNotFound).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrKeyNotFound, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)

	// Test case 6: Key is inactive
	keyInactive := &entity.Key{ID: keyID, IsActive: false}
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(keyInactive, nil).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrKeyInactive, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)

	// Test case 7: Key already reserved
	activeKeyReservation := &entity.Reservation{KeyID: keyID, Status: entity.ReservationStatusActive}
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, keyID).Return(activeKeyReservation, nil).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrKeyReserved, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 8: User already has active reservation
	userActiveReservation := &entity.Reservation{UserID: userID, Status: entity.ReservationStatusActive}
	reservation = entity.NewReservation(keyID, userID, dueAt)
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, keyID).Return(nil, entity.ErrReservationNotFound).Once()
	mockReservationRepo.On("GetByUserID", ctx, userID).Return([]*entity.Reservation{userActiveReservation}, nil).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, entity.ErrReservationAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 9: Repository error during Create
	reservation = entity.NewReservation(keyID, userID, dueAt)
	repoError := errors.New("database error")
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockKeyRepo.On("GetByID", ctx, keyID).Return(key, nil).Once()
	mockReservationRepo.On("GetActiveReservationByKey", ctx, keyID).Return(nil, entity.ErrReservationNotFound).Once()
	mockReservationRepo.On("GetByUserID", ctx, userID).Return([]*entity.Reservation{}, nil).Once()
	mockReservationRepo.On("Create", ctx, mock.AnythingOfType("*entity.Reservation")).Return(repoError).Once()

	err = reservationUseCase.CreateReservation(ctx, reservation)
	assert.Equal(t, repoError, err)
	mockUserRepo.AssertExpectations(t)
	mockKeyRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)
}

func TestReservationUseCase_ReturnKey(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()
	reservationID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	keyID := primitive.NewObjectID()

	// Test case 1: Successful return by user
	reservation := &entity.Reservation{ID: reservationID, UserID: userID, KeyID: keyID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	user := &entity.User{ID: userID, Role: entity.UserRoleResident}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err := reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Nil(t, err)
	assert.Equal(t, entity.ReservationStatusReturned, reservation.Status)
	assert.NotNil(t, reservation.ReturnedAt)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 2: Successful return by admin
	reservation = &entity.Reservation{ID: reservationID, UserID: primitive.NewObjectID(), KeyID: keyID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	adminUser := &entity.User{ID: userID, Role: entity.UserRoleAdmin}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(adminUser, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Nil(t, err)
	assert.Equal(t, entity.ReservationStatusReturned, reservation.Status)
	assert.NotNil(t, reservation.ReturnedAt)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Reservation not found
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(nil, entity.ErrReservationNotFound).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Equal(t, entity.ErrReservationNotFound, err)
	mockReservationRepo.AssertExpectations(t)

	// Test case 4: User not found
	reservation = &entity.Reservation{ID: reservationID, UserID: userID, KeyID: keyID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, entity.ErrUserNotFound).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Equal(t, entity.ErrUserNotFound, err)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Unauthorized user (not owner and not admin)
	reservation = &entity.Reservation{ID: reservationID, UserID: primitive.NewObjectID(), KeyID: keyID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	otherUser := &entity.User{ID: userID, Role: entity.UserRoleResident}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(otherUser, nil).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Equal(t, entity.ErrUnauthorized, err)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 6: Reservation not active
	reservation = &entity.Reservation{ID: reservationID, UserID: userID, KeyID: keyID, Status: entity.ReservationStatusReturned, DueAt: time.Now().Add(time.Hour)}
	user = &entity.User{ID: userID, Role: entity.UserRoleResident}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.EqualError(t, err, "reservation is not active")
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 7: Overdue reservation, user gets blocked
	reservation = &entity.Reservation{ID: reservationID, UserID: userID, KeyID: keyID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(-time.Hour)}
	user = &entity.User{ID: userID, Role: entity.UserRoleResident}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("BlockUser", ctx, userID).Return(nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err = reservationUseCase.ReturnKey(ctx, reservationID, userID)
	assert.Nil(t, err)
	assert.Equal(t, entity.ReservationStatusReturned, reservation.Status)
	assert.NotNil(t, reservation.ReturnedAt)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestReservationUseCase_GetUserReservations(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// Test case 1: Successful retrieval
	expectedReservations := []*entity.Reservation{{UserID: userID}, {UserID: userID}}
	mockReservationRepo.On("GetByUserID", ctx, userID).Return(expectedReservations, nil).Once()

	reservations, err := reservationUseCase.GetUserReservations(ctx, userID)
	assert.Nil(t, err)
	assert.Equal(t, expectedReservations, reservations)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Repository error
	repoError := errors.New("database error")
	mockReservationRepo.On("GetByUserID", ctx, userID).Return(nil, repoError).Once()

	reservations, err = reservationUseCase.GetUserReservations(ctx, userID)
	assert.Equal(t, repoError, err)
	assert.Nil(t, reservations)
	mockReservationRepo.AssertExpectations(t)
}

func TestReservationUseCase_GetAllReservations(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()

	// Test case 1: Admin user, successful retrieval
	expectedReservations := []*entity.Reservation{{}, {}}
	mockReservationRepo.On("GetAll", ctx).Return(expectedReservations, nil).Once()

	reservations, err := reservationUseCase.GetAllReservations(ctx, entity.UserRoleAdmin)
	assert.Nil(t, err)
	assert.Equal(t, expectedReservations, reservations)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	reservations, err = reservationUseCase.GetAllReservations(ctx, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: Repository error
	repoError := errors.New("database error")
	mockReservationRepo.On("GetAll", ctx).Return(nil, repoError).Once()

	reservations, err = reservationUseCase.GetAllReservations(ctx, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	assert.Nil(t, reservations)
	mockReservationRepo.AssertExpectations(t)
}

func TestReservationUseCase_ExtendReservation(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()
	reservationID := primitive.NewObjectID()
	newDueTime := time.Now().Add(24 * time.Hour)

	// Test case 1: Admin user, successful extension
	reservation := &entity.Reservation{ID: reservationID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err := reservationUseCase.ExtendReservation(ctx, reservationID, newDueTime, entity.UserRoleAdmin)
	assert.Nil(t, err)
	assert.Equal(t, newDueTime.Format(time.RFC3339), reservation.DueAt.Format(time.RFC3339))
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	err = reservationUseCase.ExtendReservation(ctx, reservationID, newDueTime, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: Reservation not found
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(nil, entity.ErrReservationNotFound).Once()

	err = reservationUseCase.ExtendReservation(ctx, reservationID, newDueTime, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrReservationNotFound, err)
	mockReservationRepo.AssertExpectations(t)

	// Test case 4: Reservation not active
	reservationInactive := &entity.Reservation{ID: reservationID, Status: entity.ReservationStatusReturned}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservationInactive, nil).Once()

	err = reservationUseCase.ExtendReservation(ctx, reservationID, newDueTime, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrCannotExtendReservation, err)
	mockReservationRepo.AssertExpectations(t)

	// Test case 5: New due time is before current due time
	reservation = &entity.Reservation{ID: reservationID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(2 * time.Hour)}
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()

	err = reservationUseCase.ExtendReservation(ctx, reservationID, time.Now().Add(time.Hour), entity.UserRoleAdmin)
	assert.EqualError(t, err, "new due time cannot be before current due time")
	mockReservationRepo.AssertExpectations(t)

	// Test case 6: Repository error during Update
	reservation = &entity.Reservation{ID: reservationID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(time.Hour)}
	repoError := errors.New("database error")
	mockReservationRepo.On("GetByID", ctx, reservationID).Return(reservation, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(repoError).Once()

	err = reservationUseCase.ExtendReservation(ctx, reservationID, newDueTime, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockReservationRepo.AssertExpectations(t)
}

func TestReservationUseCase_ProcessOverdueReservations(t *testing.T) {
	mockReservationRepo := new(mocks.ReservationRepository)
	mockKeyRepo := new(mocks.KeyRepository)
	mockUserRepo := new(mocks.UserRepository)
	reservationUseCase := usecase.NewReservationUseCase(mockReservationRepo, mockKeyRepo, mockUserRepo)

	ctx := context.Background()

	// Test case 1: No overdue reservations
	mockReservationRepo.On("GetOverdueReservations", ctx).Return([]*entity.Reservation{}, nil).Once()

	err := reservationUseCase.ProcessOverdueReservations(ctx)
	assert.Nil(t, err)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: One overdue reservation, successful processing
	userID := primitive.NewObjectID()
	reservation := &entity.Reservation{ID: primitive.NewObjectID(), UserID: userID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(-time.Hour)}
	mockReservationRepo.On("GetOverdueReservations", ctx).Return([]*entity.Reservation{reservation}, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()
	mockUserRepo.On("BlockUser", ctx, userID).Return(nil).Once()

	err = reservationUseCase.ProcessOverdueReservations(ctx)
	assert.Nil(t, err)
	assert.Equal(t, entity.ReservationStatusOverdue, reservation.Status)
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Repository error during GetOverdueReservations
	repoError := errors.New("database error")
	mockReservationRepo.On("GetOverdueReservations", ctx).Return(nil, repoError).Once()

	err = reservationUseCase.ProcessOverdueReservations(ctx)
	assert.EqualError(t, err, "failed to get overdue reservations: database error")
	mockReservationRepo.AssertExpectations(t)

	// Test case 4: Repository error during Update
	reservation = &entity.Reservation{ID: primitive.NewObjectID(), UserID: userID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(-time.Hour)}
	mockReservationRepo.On("GetOverdueReservations", ctx).Return([]*entity.Reservation{reservation}, nil).Once()
	repoError = errors.New("update error")
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(repoError).Once()

	err = reservationUseCase.ProcessOverdueReservations(ctx)
	assert.EqualError(t, err, "failed to update reservation "+reservation.ID.Hex()+": update error")
	mockReservationRepo.AssertExpectations(t)

	// Test case 5: Repository error during BlockUser
	reservation = &entity.Reservation{ID: primitive.NewObjectID(), UserID: userID, Status: entity.ReservationStatusActive, DueAt: time.Now().Add(-time.Hour)}
	mockReservationRepo.On("GetOverdueReservations", ctx).Return([]*entity.Reservation{reservation}, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()
	repoError = errors.New("block error")
	mockUserRepo.On("BlockUser", ctx, userID).Return(repoError).Once()

	err = reservationUseCase.ProcessOverdueReservations(ctx)
	assert.EqualError(t, err, "failed to block user "+userID.Hex()+": block error")
	mockReservationRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
