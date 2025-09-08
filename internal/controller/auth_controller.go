package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portaria-keys/internal/entity"
	"github.com/portaria-keys/internal/usecase"
)

// AuthController manipula requisições de autenticação
type AuthController struct {
	userUseCase *usecase.UserUseCase
}

// NewAuthController cria uma nova instância do AuthController
func NewAuthController(userUseCase *usecase.UserUseCase) *AuthController {
	return &AuthController{
		userUseCase: userUseCase,
	}
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest representa os dados de registro
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// TODO: Implementar handlers de autenticação:

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 500 {object} APIError
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}

	user, token, err := c.userUseCase.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		if err == entity.ErrInvalidCredentials || err == entity.ErrUserBlocked {
			ctx.JSON(http.StatusUnauthorized, APIError{Code: http.StatusUnauthorized, Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to login", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Register godoc
// @Summary Registro de usuário
// @Description Registra um novo usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User data"
// @Success 201 {object} entity.User
// @Failure 400 {object} APIError
// @Failure 409 {object} APIError
// @Failure 500 {object} APIError
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}

	user := entity.NewUser(req.Name, req.Email, req.Password, entity.UserRoleResident) // Default to resident

	if err := c.userUseCase.RegisterUser(ctx, user); err != nil {
		if err == entity.ErrUserAlreadyExists {
			ctx.JSON(http.StatusConflict, APIError{Code: http.StatusConflict, Message: "User with this email already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to register user", Details: err.Error()})
		return
	}

	// Avoid returning password hash
	user.Password = ""
	ctx.JSON(http.StatusCreated, user)
}