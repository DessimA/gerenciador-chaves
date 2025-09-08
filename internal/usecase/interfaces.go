package usecase

import (
	"context"
	"github.com/portaria-keys/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyRepository define as operações de persistência para chaves
type KeyRepository interface {
	Create(ctx context.Context, key *entity.Key) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Key, error)
	GetByName(ctx context.Context, name string) (*entity.Key, error)
	GetAll(ctx context.Context) ([]*entity.Key, error)
	Update(ctx context.Context, key *entity.Key) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetAvailableKeys(ctx context.Context) ([]*entity.Key, error)
}

// UserRepository define as operações de persistência para usuários
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	BlockUser(ctx context.Context, id primitive.ObjectID) error
	UnblockUser(ctx context.Context, id primitive.ObjectID) error
}

// ReservationRepository define as operações de persistência para reservas
type ReservationRepository interface {
	Create(ctx context.Context, reservation *entity.Reservation) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Reservation, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*entity.Reservation, error)
	GetByKeyID(ctx context.Context, keyID primitive.ObjectID) (*entity.Reservation, error)
	GetAll(ctx context.Context) ([]*entity.Reservation, error)
	Update(ctx context.Context, reservation *entity.Reservation) error
	GetOverdueReservations(ctx context.Context) ([]*entity.Reservation, error)
	GetActiveReservationByKey(ctx context.Context, keyID primitive.ObjectID) (*entity.Reservation, error)
}