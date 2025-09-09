package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dessima/gerenciador-chaves-api/entity"
)

// ReservationRepositoryImpl implementa a interface ReservationRepository para MongoDB.
type ReservationRepositoryImpl struct {
	collection *mongo.Collection
	tm         *TransactionManager
}

// NewReservationRepository cria uma nova instância de ReservationRepositoryImpl.
func NewReservationRepository(db *mongo.Database, client *mongo.Client) *ReservationRepositoryImpl {
	return &ReservationRepositoryImpl{
		collection: db.Collection("reservations"),
		tm:         NewTransactionManager(client),
	}
}

// Create insere uma nova reserva no banco de dados.
func (r *ReservationRepositoryImpl) Create(ctx context.Context, reservation *entity.Reservation) error {
	return r.tm.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		// Verifica se a chave já está reservada
		exists, err := r.collection.CountDocuments(sessCtx, bson.M{
			"key_id": reservation.KeyID,
			"status": entity.ReservationStatusActive,
		})
		if err != nil {
			return err
		}
		if exists > 0 {
			return entity.ErrKeyAlreadyReserved
		}

		// Inicializa a versão e timestamps
		reservation.Version = 1
		reservation.CreatedAt = time.Now()
		reservation.UpdatedAt = reservation.CreatedAt

		// Insere a reserva
		_, err = r.collection.InsertOne(sessCtx, reservation)
		return err
	})
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

// Update atualiza uma reserva existente usando bloqueio otimista
func (r *ReservationRepositoryImpl) Update(ctx context.Context, reservation *entity.Reservation) error {
	return r.tm.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		// Incrementa a versão antes da atualização
		currentVersion := reservation.Version
		reservation.Version++
		reservation.UpdatedAt = time.Now()

		result, err := r.collection.UpdateOne(
			sessCtx,
			bson.M{
				"_id":     reservation.ID,
				"version": currentVersion,
			},
			bson.M{"$set": reservation},
		)

		if err != nil {
			return err
		}

		if result.ModifiedCount == 0 {
			return ErrConcurrentModification
		}

		return nil
	})
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

// GetActiveReservationByKey busca uma reserva ativa por KeyID.
func (r *ReservationRepositoryImpl) GetActiveReservationByKey(ctx context.Context, keyID primitive.ObjectID) (*entity.Reservation, error) {
	var reservation entity.Reservation
	err := r.collection.FindOne(ctx, bson.M{
		"key_id": keyID,
		"status": entity.ReservationStatusActive,
	}).Decode(&reservation)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

// GetOverdueReservations retorna todas as reservas vencidas.
func (r *ReservationRepositoryImpl) GetOverdueReservations(ctx context.Context) ([]*entity.Reservation, error) {
	var reservations []*entity.Reservation
	cursor, err := r.collection.Find(ctx, bson.M{
		"status": entity.ReservationStatusActive,
		"return_date": bson.M{
			"$lt": time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}
