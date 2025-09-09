package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/usecase"
	"github.com/dessima/gerenciador-chaves-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserUseCase_RegisterUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockReservationRepo := new(mocks.ReservationRepository) // Not used in RegisterUser
	userUseCase := usecase.NewUserUseCase(mockUserRepo, mockReservationRepo)

	ctx := context.Background()

	// Test case 1: Successful registration
	user := entity.NewUser("Test User", "test@example.com", "password123", entity.UserRoleResident)
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(nil, entity.ErrUserNotFound).Once()
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil).Once()

	err := userUseCase.RegisterUser(ctx, user)
	assert.Nil(t, err)
	assert.NotEqual(t, "password123", user.Password) // Password should be hashed
	mockUserRepo.AssertExpectations(t)

	// Test case 2: Invalid user data
	user = entity.NewUser("", "invalid-email", "short", entity.UserRoleResident)
	err = userUseCase.RegisterUser(ctx, user)
	assert.NotNil(t, err)

	// Test case 3: User with same email already exists
	user = entity.NewUser("Existing User", "existing@example.com", "password123", entity.UserRoleResident)
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()

	err = userUseCase.RegisterUser(ctx, user)
	assert.Equal(t, entity.ErrUserAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 4: Repository error during GetByEmail
	user = entity.NewUser("Repo Error User", "repo@example.com", "password123", entity.UserRoleResident)
	repoError := errors.New("database error")
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(nil, repoError).Once()

	err = userUseCase.RegisterUser(ctx, user)
	assert.Equal(t, repoError, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Repository error during Create
	user = entity.NewUser("Create Error User", "create@example.com", "password123", entity.UserRoleResident)
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(nil, entity.ErrUserNotFound).Once()
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(repoError).Once()

	err = userUseCase.RegisterUser(ctx, user)
	assert.Equal(t, repoError, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUseCase_LoginUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	userUseCase := usecase.NewUserUseCase(mockUserRepo, mockReservationRepo)

	ctx := context.Background()

	// Test case 1: Successful login
	user := entity.NewUser("Test User", "test@example.com", "password123", entity.UserRoleResident)
	user.HashPassword() // Hash the password before returning
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()

	loggedInUser, token, err := userUseCase.LoginUser(ctx, user.Email, "password123")
	assert.Nil(t, err)
	assert.Equal(t, user.Email, loggedInUser.Email)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)

	// Test case 2: User not found
	mockUserRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, entity.ErrUserNotFound).Once()

	loggedInUser, token, err = userUseCase.LoginUser(ctx, "nonexistent@example.com", "password123")
	assert.Equal(t, entity.ErrInvalidCredentials, err)
	assert.Nil(t, loggedInUser)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Incorrect password
	user = entity.NewUser("Test User", "wrongpass@example.com", "password123", entity.UserRoleResident)
	user.HashPassword()
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()

	loggedInUser, token, err = userUseCase.LoginUser(ctx, user.Email, "wrong_password")
	assert.Equal(t, entity.ErrInvalidCredentials, err)
	assert.Nil(t, loggedInUser)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)

	// Test case 4: User is blocked
	user = entity.NewUser("Blocked User", "blocked@example.com", "password123", entity.UserRoleResident)
	user.HashPassword()
	user.IsBlocked = true
	mockUserRepo.On("GetByEmail", ctx, user.Email).Return(user, nil).Once()

	loggedInUser, token, err = userUseCase.LoginUser(ctx, user.Email, "password123")
	assert.Equal(t, entity.ErrUserBlocked, err)
	assert.Nil(t, loggedInUser)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Repository error during GetByEmail
	repoError := errors.New("database error")
	mockUserRepo.On("GetByEmail", ctx, "repoerror@example.com").Return(nil, repoError).Once()

	loggedInUser, token, err = userUseCase.LoginUser(ctx, "repoerror@example.com", "password123")
	assert.Equal(t, entity.ErrInvalidCredentials, err)
	assert.Nil(t, loggedInUser)
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestUserUseCase_BlockUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	userUseCase := usecase.NewUserUseCase(mockUserRepo, mockReservationRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// Test case 1: Admin user, successful block
	user := &entity.User{ID: userID, Name: "Test User", Email: "test@example.com", Role: entity.UserRoleResident, IsBlocked: false}
	reservation := &entity.Reservation{ID: primitive.NewObjectID(), UserID: userID, Status: entity.ReservationStatusActive}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("BlockUser", ctx, userID).Return(nil).Once()
	mockReservationRepo.On("GetByUserID", ctx, userID).Return([]*entity.Reservation{reservation}, nil).Once()
	mockReservationRepo.On("Update", ctx, mock.AnythingOfType("*entity.Reservation")).Return(nil).Once()

	err := userUseCase.BlockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Nil(t, err)
	mockUserRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	err = userUseCase.BlockUser(ctx, userID, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: User not found
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, entity.ErrUserNotFound).Once()

	err = userUseCase.BlockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 4: User already blocked
	userBlocked := &entity.User{ID: userID, Name: "Blocked User", Email: "blocked@example.com", Role: entity.UserRoleResident, IsBlocked: true}
	mockUserRepo.On("GetByID", ctx, userID).Return(userBlocked, nil).Once()

	err = userUseCase.BlockUser(ctx, userID, entity.UserRoleAdmin)
	assert.EqualError(t, err, "user is already blocked")
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Repository error during BlockUser
	user = &entity.User{ID: userID, Name: "Test User", Email: "test@example.com", Role: entity.UserRoleResident, IsBlocked: false}
	repoError := errors.New("database error")
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("BlockUser", ctx, userID).Return(repoError).Once()

	err = userUseCase.BlockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 6: Repository error during GetByUserID (reservations)
	user = &entity.User{ID: userID, Name: "Test User", Email: "test@example.com", Role: entity.UserRoleResident, IsBlocked: false}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("BlockUser", ctx, userID).Return(nil).Once()
	mockReservationRepo.On("GetByUserID", ctx, userID).Return(nil, repoError).Once()

	err = userUseCase.BlockUser(ctx, userID, entity.UserRoleAdmin)
	assert.EqualError(t, err, "failed to get user reservations: database error")
	mockUserRepo.AssertExpectations(t)
	mockReservationRepo.AssertExpectations(t)
}

func TestUserUseCase_UnblockUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockReservationRepo := new(mocks.ReservationRepository)
	userUseCase := usecase.NewUserUseCase(mockUserRepo, mockReservationRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// Test case 1: Admin user, successful unblock
	user := &entity.User{ID: userID, Name: "Test User", Email: "test@example.com", Role: entity.UserRoleResident, IsBlocked: true}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("UnblockUser", ctx, userID).Return(nil).Once()

	err := userUseCase.UnblockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Nil(t, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 2: Non-admin user, unauthorized
	err = userUseCase.UnblockUser(ctx, userID, entity.UserRoleResident)
	assert.Equal(t, entity.ErrUnauthorized, err)

	// Test case 3: User not found
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, entity.ErrUserNotFound).Once()

	err = userUseCase.UnblockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Equal(t, entity.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)

	// Test case 4: User already unblocked
	userUnblocked := &entity.User{ID: userID, Name: "Unblocked User", Email: "unblocked@example.com", Role: entity.UserRoleResident, IsBlocked: false}
	mockUserRepo.On("GetByID", ctx, userID).Return(userUnblocked, nil).Once()

	err = userUseCase.UnblockUser(ctx, userID, entity.UserRoleAdmin)
	assert.EqualError(t, err, "user is not blocked")
	mockUserRepo.AssertExpectations(t)

	// Test case 5: Repository error during UnblockUser
	user = &entity.User{ID: userID, Name: "Test User", Email: "test@example.com", Role: entity.UserRoleResident, IsBlocked: true}
	repoError := errors.New("database error")
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
	mockUserRepo.On("UnblockUser", ctx, userID).Return(repoError).Once()

	err = userUseCase.UnblockUser(ctx, userID, entity.UserRoleAdmin)
	assert.Equal(t, repoError, err)
	mockUserRepo.AssertExpectations(t)
}
