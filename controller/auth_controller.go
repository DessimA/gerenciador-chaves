package controller

import (
	"net/http"

	"github.com/dessima/gerenciador-chaves-api/controller/dto"
	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/http/response"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/validation"
	"github.com/dessima/gerenciador-chaves-api/usecase"
	"github.com/gin-gonic/gin"
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

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequest true "Credenciais do usuário"
// @Success 200 {object} response.APIResponse{data=dto.AuthResponse} "Login bem sucedido"
// @Failure 400 {object} response.APIResponse{error=response.APIError} "Erro de validação"
// @Failure 401 {object} response.APIResponse{error=response.APIError} "Credenciais inválidas"
// @Failure 500 {object} response.APIResponse{error=response.APIError} "Erro interno"
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq dto.UserLoginRequest
	if err := validation.ValidateRequest(ctx, &loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user, token, err := c.userUseCase.LoginUser(ctx.Request.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		var statusCode int
		switch err {
		case entity.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, response.NewError(statusCode, err.Error()))
		return
	}

	authResponse := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			IsBlocked: user.IsBlocked,
		},
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(authResponse, nil))
}

// Register godoc
// @Summary Registra novo usuário
// @Description Cria um novo usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "Dados do usuário"
// @Success 201 {object} response.APIResponse{data=dto.UserResponse} "Usuário registrado com sucesso"
// @Failure 400 {object} response.APIResponse{error=response.APIError} "Erro de validação"
// @Failure 409 {object} response.APIResponse{error=response.APIError} "Email já cadastrado"
// @Failure 500 {object} response.APIResponse{error=response.APIError} "Erro interno"
// @Router /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var registerReq dto.UserRegisterRequest
	if err := validation.ValidateRequest(ctx, &registerReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	user := &entity.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: registerReq.Password,
		Role:     entity.UserRoleResident, // Por padrão, novos usuários são residentes
	}

	if err := c.userUseCase.RegisterUser(ctx.Request.Context(), user); err != nil {
		var statusCode int
		switch err {
		case entity.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		ctx.JSON(statusCode, response.NewError(statusCode, err.Error()))
		return
	}

	userResponse := dto.UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		IsBlocked: user.IsBlocked,
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(userResponse, nil))
}
