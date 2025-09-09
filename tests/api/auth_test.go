package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/dessima/gerenciador-chaves-api/controller"
	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/repository"
	"github.com/dessima/gerenciador-chaves-api/usecase"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var router *gin.Engine
var testDB *mongo.Database

func setupAPITest(t *testing.T) {
	// Connect to a test MongoDB instance
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.Nil(t, err)
	testDB = client.Database("test_api_db_" + primitive.NewObjectID().Hex())

	// Drop the database before each test to ensure a clean state
	_ = testDB.Drop(context.Background())

	// Initialize Repositories
	userRepo := repository.NewUserRepository(testDB)
	keyRepo := repository.NewKeyRepository(testDB)
	reservationRepo := repository.NewReservationRepository(testDB)

	// Initialize Use Cases
	userUseCase := usecase.NewUserUseCase(userRepo, reservationRepo)
	keyUseCase := usecase.NewKeyUseCase(keyRepo, reservationRepo)
	reservationUseCase := usecase.NewReservationUseCase(reservationRepo, keyRepo, userRepo)

	// Setup Router
	router = gin.Default()
	router.POST("/auth/register", controller.NewAuthController(userUseCase).Register)
	router.POST("/auth/login", controller.NewAuthController(userUseCase).Login)
}

func teardownAPITest() {
	_ = testDB.Drop(context.Background())
	_ = testDB.Client().Disconnect(context.Background())
}

func TestAuth_Register(t *testing.T) {
	setupAPITest(t)
	defer teardownAPITest()

	// Test case 1: Successful registration
	registerReq := controller.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var user entity.User
	_ = json.Unmarshal(rec.Body.Bytes(), &user)
	assert.Equal(t, registerReq.Email, user.Email)
	assert.Empty(t, user.Password) // Password should not be returned

	// Test case 2: Registration with existing email
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	var apiError controller.APIError
	_ = json.Unmarshal(rec.Body.Bytes(), &apiError)
	assert.Equal(t, "User with this email already exists", apiError.Message)

	// Test case 3: Invalid request body
	invalidBody := []byte(`{"email": "invalid"}`)
	req, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuth_Login(t *testing.T) {
	setupAPITest(t)
	defer teardownAPITest()

	// Register a user first
	registerReq := controller.RegisterRequest{
		Name:     "Login User",
		Email:    "login@example.com",
		Password: "loginpassword",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Test case 1: Successful login
	loginReq := controller.LoginRequest{
		Email:    "login@example.com",
		Password: "loginpassword",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")

	// Test case 2: Invalid credentials
	loginReq.Password = "wrongpassword"
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	_ = json.Unmarshal(rec.Body.Bytes(), &apiError)
	assert.Equal(t, "credenciais inválidas", apiError.Message)

	// Test case 3: User not found
	loginReq.Email = "nonexistent@example.com"
	loginReq.Password = "anypassword"
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	_ = json.Unmarshal(rec.Body.Bytes(), &apiError)
	assert.Equal(t, "credenciais inválidas", apiError.Message)
}
