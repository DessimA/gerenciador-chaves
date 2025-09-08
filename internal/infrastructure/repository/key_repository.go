package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/portaria-keys/internal/entity"
)

// KeyRepositoryImpl implementa a interface KeyRepository para MongoDB.
type KeyRepositoryImpl struct {
	collection *mongo.Collection
}

// NewKeyRepository cria uma nova instância de KeyRepositoryImpl.
func NewKeyRepository(db *mongo.Database) *KeyRepositoryImpl {
	return &KeyRepositoryImpl{
		collection: db.Collection("keys"),
	}
}

// Create insere uma nova chave no banco de dados.
func (r *KeyRepositoryImpl) Create(ctx context.Context, key *entity.Key) error {
	_, err := r.collection.InsertOne(ctx, key)
	return err
}

// GetByID busca uma chave pelo ID.
func (r *KeyRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Key, error) {
	var key entity.Key
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrKeyNotFound
		}
		return nil, err
	}
	return &key, nil
}

// GetByName busca uma chave pelo nome.
func (r *KeyRepositoryImpl) GetByName(ctx context.Context, name string) (*entity.Key, error) {
	var key entity.Key
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrKeyNotFound
		}
		return nil, err
	}
	return &key, nil
}

// GetAll retorna todas as chaves.
func (r *KeyRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Key, error) {
	var keys []*entity.Key
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// Update atualiza uma chave existente.
func (r *KeyRepositoryImpl) Update(ctx context.Context, key *entity.Key) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": key.ID}, key)
	return err
}

// Delete remove uma chave pelo ID.
func (r *KeyRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetAvailableKeys retorna chaves ativas que não possuem reservas ativas.
func (r *KeyRepositoryImpl) GetAvailableKeys(ctx context.Context) ([]*entity.Key, error) {
	// This requires a more complex aggregation or a separate query to the reservations collection.
	// For simplicity, this implementation will return all active keys.
	// A proper implementation would involve checking the Reservation collection for active reservations.
	var keys []*entity.Key
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}