package repository_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContainer testcontainers.Container
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()

	// Setup MongoDB container
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start mongo container: %s", err)
	}
	mongoContainer = mongoC

	endpoint, err := mongoContainer.Endpoint(ctx, "mongodb")
	if err != nil {
		log.Fatalf("Could not get mongo endpoint: %s", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("Could not connect to mongo: %s", err)
	}
	mongoClient = client

	// Run tests
	code := m.Run()

	// Teardown
	if err := mongoContainer.Terminate(ctx); err != nil {
		log.Fatalf("Could not terminate mongo container: %s", err)
	}
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Could not disconnect from mongo: %s", err)
	}
	os.Exit(code)
}

func setupTest(t *testing.T) *repository.KeyRepositoryImpl {
	// Use a unique database for each test to ensure isolation
	dbName := "test_db_" + primitive.NewObjectID().Hex()
	db := mongoClient.Database(dbName)

	// Drop the database before each test to ensure a clean state
	if err := db.Drop(ctx); err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}

	return repository.NewKeyRepository(db)
}

func TestKeyRepository_Create(t *testing.T) {
	repo := setupTest(t)
	key := entity.NewKey("Test Key", "Description")

	err := repo.Create(ctx, key)
	assert.Nil(t, err)
	assert.False(t, key.ID.IsZero())

	foundKey, err := repo.GetByID(ctx, key.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, key.Name, foundKey.Name)
}

func TestKeyRepository_GetByID(t *testing.T) {
	repo := setupTest(t)
	key := entity.NewKey("Test Key", "Description")
	repo.Create(ctx, key)

	foundKey, err := repo.GetByID(ctx, key.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, key.Name, foundKey.Name)

	// Test not found
	_, err = repo.GetByID(ctx, primitive.NewObjectID())
	assert.Equal(t, entity.ErrKeyNotFound, err)
}

func TestKeyRepository_GetByName(t *testing.T) {
	repo := setupTest(t)
	key := entity.NewKey("Unique Key Name", "Description")
	repo.Create(ctx, key)

	foundKey, err := repo.GetByName(ctx, key.Name)
	assert.Nil(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, key.Name, foundKey.Name)

	// Test not found
	_, err = repo.GetByName(ctx, "NonExistentKey")
	assert.Equal(t, entity.ErrKeyNotFound, err)
}

func TestKeyRepository_GetAll(t *testing.T) {
	repo := setupTest(t)
	repo.Create(ctx, entity.NewKey("Key1", "Desc1"))
	repo.Create(ctx, entity.NewKey("Key2", "Desc2"))

	keys, err := repo.GetAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, keys, 2)
}

func TestKeyRepository_Update(t *testing.T) {
	repo := setupTest(t)
	key := entity.NewKey("Original Key", "Original Desc")
	repo.Create(ctx, key)

	key.Name = "Updated Key"
	key.Description = "Updated Desc"
	key.IsActive = false
	key.UpdatedAt = time.Now()

	err := repo.Update(ctx, key)
	assert.Nil(t, err)

	foundKey, err := repo.GetByID(ctx, key.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, "Updated Key", foundKey.Name)
	assert.Equal(t, "Updated Desc", foundKey.Description)
	assert.False(t, foundKey.IsActive)
}

func TestKeyRepository_Delete(t *testing.T) {
	repo := setupTest(t)
	key := entity.NewKey("Key To Delete", "Desc")
	repo.Create(ctx, key)

	_, err := repo.GetByID(ctx, key.ID)
	assert.Nil(t, err)

	err = repo.Delete(ctx, key.ID)
	assert.Nil(t, err)

	_, err = repo.GetByID(ctx, key.ID)
	assert.Equal(t, entity.ErrKeyNotFound, err)
}

func TestKeyRepository_GetAvailableKeys(t *testing.T) {
	repo := setupTest(t)
	activeKey := entity.NewKey("Active Key", "Desc")
	activeKey.IsActive = true
	repo.Create(ctx, activeKey)

	inactiveKey := entity.NewKey("Inactive Key", "Desc")
	inactiveKey.IsActive = false
	repo.Create(ctx, inactiveKey)

	// Note: This test currently only checks for IsActive=true.
	// A full implementation would also check the Reservation collection for active reservations.
	keys, err := repo.GetAvailableKeys(ctx)
	assert.Nil(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, activeKey.Name, keys[0].Name)
}
