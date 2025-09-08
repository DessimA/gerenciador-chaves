package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/portaria-keys/internal/entity"
	"github.com/portaria-keys/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var reservationMongoClient *mongo.Client
var reservationMongoContainer testcontainers.Container
var reservationCtx context.Context

func TestMainReservation(m *testing.M) {
	reservationCtx = context.Background()

	// Setup MongoDB container
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(reservationCtx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start mongo container for reservation tests: %s", err)
	}
	reservationMongoContainer = mongoC

	endpoint, err := reservationMongoContainer.Endpoint(reservationCtx, "mongodb")
	if err != nil {
		log.Fatalf("Could not get mongo endpoint for reservation tests: %s", err)
	}

	client, err := mongo.Connect(reservationCtx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("Could not connect to mongo for reservation tests: %s", err)
	}
	reservationMongoClient = client

	// Run tests
	code := m.Run()

	// Teardown
	if err := reservationMongoContainer.Terminate(reservationCtx); err != nil {
		log.Fatalf("Could not terminate mongo container for reservation tests: %s", err)
	}
	if err := reservationMongoClient.Disconnect(reservationCtx); err != nil {
		log.Fatalf("Could not disconnect from mongo for reservation tests: %s", err)
	}
	os.Exit(code)
}

func setupReservationTest(t *testing.T) *repository.ReservationRepositoryImpl {
	// Use a unique database for each test to ensure isolation
	dbName := "test_reservation_db_" + primitive.NewObjectID().Hex()
	db := reservationMongoClient.Database(dbName)

	// Drop the database before each test to ensure a clean state
	if err := db.Drop(reservationCtx); err != nil {
		log.Fatalf("Failed to drop test reservation database: %v", err)
	}

	return repository.NewReservationRepository(db)
}

func TestReservationRepository_Create(t *testing.T) {
	repo := setupReservationTest(t)
	reservation := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour))

	err := repo.Create(reservationCtx, reservation)
	assert.Nil(t, err)
	assert.False(t, reservation.ID.IsZero())

	foundReservation, err := repo.GetByID(reservationCtx, reservation.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundReservation)
	assert.Equal(t, reservation.KeyID, foundReservation.KeyID)
}

func TestReservationRepository_GetByID(t *testing.T) {
	repo := setupReservationTest(t)
	reservation := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, reservation)

	foundReservation, err := repo.GetByID(reservationCtx, reservation.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundReservation)
	assert.Equal(t, reservation.KeyID, foundReservation.KeyID)

	// Test not found
	_, err = repo.GetByID(reservationCtx, primitive.NewObjectID())
	assert.Equal(t, entity.ErrReservationNotFound, err)
}

func TestReservationRepository_GetByUserID(t *testing.T) {
	repo := setupReservationTest(t)
	userID := primitive.NewObjectID()
	reservation1 := entity.NewReservation(primitive.NewObjectID(), userID, time.Now().Add(time.Hour))
	reservation2 := entity.NewReservation(primitive.NewObjectID(), userID, time.Now().Add(2*time.Hour))
	reservation3 := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, reservation1)
	repo.Create(reservationCtx, reservation2)
	repo.Create(reservationCtx, reservation3)

	foundReservations, err := repo.GetByUserID(reservationCtx, userID)
	assert.Nil(t, err)
	assert.Len(t, foundReservations, 2)
	assert.Contains(t, foundReservations, reservation1)
	assert.Contains(t, foundReservations, reservation2)
}

func TestReservationRepository_GetByKeyID(t *testing.T) {
	repo := setupReservationTest(t)
	keyID := primitive.NewObjectID()
	reservation := entity.NewReservation(keyID, primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, reservation)

	foundReservation, err := repo.GetByKeyID(reservationCtx, keyID)
	assert.Nil(t, err)
	assert.NotNil(t, foundReservation)
	assert.Equal(t, reservation.KeyID, foundReservation.KeyID)

	// Test not found
	_, err = repo.GetByKeyID(reservationCtx, primitive.NewObjectID())
	assert.Equal(t, entity.ErrReservationNotFound, err)
}

func TestReservationRepository_GetAll(t *testing.T) {
	repo := setupReservationTest(t)
	repo.Create(reservationCtx, entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour)))
	repo.Create(reservationCtx, entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour)))

	reservations, err := repo.GetAll(reservationCtx)
	assert.Nil(t, err)
	assert.Len(t, reservations, 2)
}

func TestReservationRepository_Update(t *testing.T) {
	repo := setupReservationTest(t)
	reservation := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, reservation)

	reservation.Status = entity.ReservationStatusReturned
	reservation.UpdatedAt = time.Now()

	err := repo.Update(reservationCtx, reservation)
	assert.Nil(t, err)

	foundReservation, err := repo.GetByID(reservationCtx, reservation.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundReservation)
	assert.Equal(t, entity.ReservationStatusReturned, foundReservation.Status)
}

func TestReservationRepository_GetOverdueReservations(t *testing.T) {
	repo := setupReservationTest(t)
	// Overdue reservation
	overdueRes := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(-time.Hour))
	repo.Create(reservationCtx, overdueRes)

	// Active reservation (not overdue)
	activeRes := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, activeRes)

	// Returned reservation
	returnedRes := entity.NewReservation(primitive.NewObjectID(), primitive.NewObjectID(), time.Now().Add(-2*time.Hour))
	returnedRes.Status = entity.ReservationStatusReturned
	repo.Create(reservationCtx, returnedRes)

	overdueReservations, err := repo.GetOverdueReservations(reservationCtx)
	assert.Nil(t, err)
	assert.Len(t, overdueReservations, 1)
	assert.Equal(t, overdueRes.ID, overdueReservations[0].ID)
}

func TestReservationRepository_GetActiveReservationByKey(t *testing.T) {
	repo := setupReservationTest(t)
	keyID := primitive.NewObjectID()
	activeRes := entity.NewReservation(keyID, primitive.NewObjectID(), time.Now().Add(time.Hour))
	repo.Create(reservationCtx, activeRes)

	// Inactive reservation for the same key
	inactiveRes := entity.NewReservation(keyID, primitive.NewObjectID(), time.Now().Add(time.Hour))
	inactiveRes.Status = entity.ReservationStatusReturned
	repo.Create(reservationCtx, inactiveRes)

	foundReservation, err := repo.GetActiveReservationByKey(reservationCtx, keyID)
	assert.Nil(t, err)
	assert.NotNil(t, foundReservation)
	assert.Equal(t, activeRes.ID, foundReservation.ID)

	// Test not found (no active reservation for another key)
	_, err = repo.GetActiveReservationByKey(reservationCtx, primitive.NewObjectID())
	assert.Equal(t, entity.ErrReservationNotFound, err)
}
