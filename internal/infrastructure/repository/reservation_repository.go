package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/portaria-keys/internal/entity"
)

// ReservationRepositoryImpl implementa a interface ReservationRepository para MongoDB.
type ReservationRepositoryImpl struct {
	collection *mongo.Collection
}

// NewReservationRepository cria uma nova instância de ReservationRepositoryImpl.
func NewReservationRepository(db *mongo.Database) *ReservationRepositoryImpl {
	return &ReservationRepositoryImpl{
		collection: db.Collection("reservations"),
	}
}

// Create insere uma nova reserva no banco de dados.
func (r *ReservationRepositoryImpl) Create(ctx context.Context, reservation *entity.Reservation) error {
	_, err := r.collection.InsertOne(ctx, reservation)
	return err
}

// GetByID busca uma reserva pelo ID.
func (r *ReservationRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&reservation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

// GetByUserID busca reservas de um usuário.
func (r *ReservationRepositoryImpl) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*entity.Reservation, error) {
	var reservations []*entity.Reservation
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetByKeyID busca uma reserva ativa por KeyID.
func (r *ReservationRepositoryImpl) GetByKeyID(ctx context.Context, keyID primitive.ObjectID) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.collection.FindOne(ctx, bson.M{"key_id": keyID, "status": entity.ReservationStatusActive}).Decode(&reservation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

// GetAll retorna todas as reservas.
func (r *ReservationRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Reservation, error) {
	var reservations []*entity.Reservation
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

// Update atualiza uma reserva existente.
func (r *ReservationRepositoryImpl) Update(ctx context.Context, reservation *entity.Reservation) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": reservation.ID}, reservation)
	return err
}

// GetOverdueReservations retorna reservas em atraso.
func (r *ReservationRepositoryImpl) GetOverdueReservations(ctx context.Context) ([]*entity.Reservation, error) {
	var reservations []*entity.Reservation
	filter := bson.M{
		"status": entity.ReservationStatusActive,
		"due_at": bson.M{"$lt": time.Now()},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetActiveReservationByKey busca uma reserva ativa por KeyID.
func (r *ReservationRepositoryImpl) GetActiveReservationByKey(ctx context.Context, keyID primitive.ObjectID) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.collection.FindOne(ctx, bson.M{"key_id": keyID, "status": entity.ReservationStatusActive}).Decode(&reservation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}