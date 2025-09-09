package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserUseCase implementa os casos de uso relacionados aos usuários
type UserUseCase struct {
	config          *config.Config
	userRepo        UserRepository
	reservationRepo ReservationRepository
}

// NewUserUseCase cria uma nova instância do UserUseCase
func NewUserUseCase(userRepo UserRepository, reservationRepo ReservationRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		reservationRepo: reservationRepo,
	}
}

// Authenticate é um alias para LoginUser que mantém a compatibilidade da interface
// Deprecated: Usar LoginUser em vez deste método
func (uc *UserUseCase) Authenticate(ctx context.Context, email, password string) (string, error) {
	_, token, err := uc.LoginUser(ctx, email, password)
	return token, err
}

// RegisterUser registra um novo usuário
func (uc *UserUseCase) RegisterUser(ctx context.Context, user *entity.User) error {
	// Verificar se o usuário já existe
	existingUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && err != entity.ErrUserNotFound {
		return err
	}
	if existingUser != nil {
		return entity.ErrUserAlreadyExists
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return uc.userRepo.Create(ctx, user)
}

// LoginUser autentica um usuário
func (uc *UserUseCase) LoginUser(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", entity.ErrInvalidCredentials
	}

	if user.IsBlocked {
		return nil, "", entity.ErrUserBlocked
	}

	if err := user.ValidatePassword(password); err != nil {
		return nil, "", entity.ErrInvalidCredentials
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token válido por 24 horas
	})

	tokenString, err := token.SignedString([]byte(uc.config.JWTSecret))
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, tokenString, nil
}

// BlockUser bloqueia um usuário (apenas admin)
func (uc *UserUseCase) BlockUser(ctx context.Context, userID primitive.ObjectID, adminRole entity.UserRole) error {
	if adminRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	if user.IsBlocked {
		return errors.New("user is already blocked")
	}

	// Block the user
	if err := uc.userRepo.BlockUser(ctx, userID); err != nil {
		return err
	}

	// Cancel active reservations for the user
	reservations, err := uc.reservationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user reservations: %w", err)
	}

	for _, res := range reservations {
		if res.Status == entity.ReservationStatusActive {
			if err := res.MarkAsOverdue(); err != nil {
				return fmt.Errorf("failed to mark reservation %s as overdue: %w", res.ID.Hex(), err)
			}
			if err := uc.reservationRepo.Update(ctx, res); err != nil {
				return fmt.Errorf("failed to update reservation %s: %w", res.ID.Hex(), err)
			}
		}
	}

	return nil
}

// UnblockUser desbloqueia um usuário (apenas admin)
func (uc *UserUseCase) UnblockUser(ctx context.Context, userID primitive.ObjectID, adminRole entity.UserRole) error {
	if adminRole != entity.UserRoleAdmin {
		return entity.ErrUnauthorized
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return entity.ErrUserNotFound
	}

	if !user.IsBlocked {
		return errors.New("user is not blocked")
	}

	return uc.userRepo.UnblockUser(ctx, userID)
}
