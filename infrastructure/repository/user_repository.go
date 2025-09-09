package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dessima/gerenciador-chaves-api/entity"
)

// UserRepositoryImpl implementa a interface UserRepository para MongoDB.
type UserRepositoryImpl struct {
	collection *mongo.Collection
}

// NewUserRepository cria uma nova instância de UserRepositoryImpl.
func NewUserRepository(db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		collection: db.Collection("users"),
	}
}

// Create insere um novo usuário no banco de dados.
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// GetByID busca um usuário pelo ID.
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail busca um usuário pelo email.
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update atualiza um usuário existente.
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// BlockUser bloqueia um usuário pelo ID.
func (r *UserRepositoryImpl) BlockUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"is_blocked": true}})
	return err
}

// UnblockUser desbloqueia um usuário pelo ID.
func (r *UserRepositoryImpl) UnblockUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"is_blocked": false}})
	return err
}
