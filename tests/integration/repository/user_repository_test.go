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

var userMongoClient *mongo.Client
var userMongoContainer testcontainers.Container
var userCtx context.Context

func TestMainUser(m *testing.M) {
	userCtx = context.Background()

	// Setup MongoDB container
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	mongoC, err := testcontainers.GenericContainer(userCtx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start mongo container for user tests: %s", err)
	}
	userMongoContainer = mongoC

	endpoint, err := userMongoContainer.Endpoint(userCtx, "mongodb")
	if err != nil {
		log.Fatalf("Could not get mongo endpoint for user tests: %s", err)
	}

	client, err := mongo.Connect(userCtx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("Could not connect to mongo for user tests: %s", err)
	}
	userMongoClient = client

	// Run tests
	code := m.Run()

	// Teardown
	if err := userMongoContainer.Terminate(userCtx); err != nil {
		log.Fatalf("Could not terminate mongo container for user tests: %s", err)
	}
	if err := userMongoClient.Disconnect(userCtx); err != nil {
		log.Fatalf("Could not disconnect from mongo for user tests: %s", err)
	}
	os.Exit(code)
}

func setupUserTest(t *testing.T) *repository.UserRepositoryImpl {
	// Use a unique database for each test to ensure isolation
	dbName := "test_user_db_" + primitive.NewObjectID().Hex()
	db := userMongoClient.Database(dbName)

	// Drop the database before each test to ensure a clean state
	if err := db.Drop(userCtx); err != nil {
		log.Fatalf("Failed to drop test user database: %v", err)
	}

	return repository.NewUserRepository(db)
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupUserTest(t)
	user := entity.NewUser("Test User", "test@example.com", "password", entity.UserRoleResident)

	err := repo.Create(userCtx, user)
	assert.Nil(t, err)
	assert.False(t, user.ID.IsZero())

	foundUser, err := repo.GetByID(userCtx, user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := setupUserTest(t)
	user := entity.NewUser("Test User", "test@example.com", "password", entity.UserRoleResident)
	repo.Create(userCtx, user)

	foundUser, err := repo.GetByID(userCtx, user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test not found
	_, err = repo.GetByID(userCtx, primitive.NewObjectID())
	assert.Equal(t, entity.ErrUserNotFound, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := setupUserTest(t)
	user := entity.NewUser("Test User", "unique@example.com", "password", entity.UserRoleResident)
	repo.Create(userCtx, user)

	foundUser, err := repo.GetByEmail(userCtx, user.Email)
	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test not found
	_, err = repo.GetByEmail(userCtx, "nonexistent@example.com")
	assert.Equal(t, entity.ErrUserNotFound, err)
}

func TestUserRepository_Update(t *testing.T) {
	repo := setupUserTest(t)
	user := entity.NewUser("Original User", "original@example.com", "password", entity.UserRoleResident)
	repo.Create(userCtx, user)

	user.Name = "Updated User"
	user.Email = "updated@example.com"
	user.UpdatedAt = time.Now()

	err := repo.Update(userCtx, user)
	assert.Nil(t, err)

	foundUser, err := repo.GetByID(userCtx, user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, "Updated User", foundUser.Name)
	assert.Equal(t, "updated@example.com", foundUser.Email)
}

func TestUserRepository_BlockAndUnblockUser(t *testing.T) {
	repo := setupUserTest(t)
	user := entity.NewUser("Blockable User", "block@example.com", "password", entity.UserRoleResident)
	repo.Create(userCtx, user)

	// Block user
	err := repo.BlockUser(userCtx, user.ID)
	assert.Nil(t, err)

	foundUser, err := repo.GetByID(userCtx, user.ID)
	assert.Nil(t, err)
	assert.True(t, foundUser.IsBlocked)

	// Unblock user
	err = repo.UnblockUser(userCtx, user.ID)
	assert.Nil(t, err)

	foundUser, err = repo.GetByID(userCtx, user.ID)
	assert.Nil(t, err)
	assert.False(t, foundUser.IsBlocked)
}
