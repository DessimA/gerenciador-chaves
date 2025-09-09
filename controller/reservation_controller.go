package controller

import (
	"net/http"
	"time"

	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/http/response"
	"github.com/dessima/gerenciador-chaves-api/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReservationController manipula requisições relacionadas às reservas
type ReservationController struct {
	reservationUseCase *usecase.ReservationUseCase
}

// NewReservationController cria uma nova instância do ReservationController
func NewReservationController(reservationUseCase *usecase.ReservationUseCase) *ReservationController {
	return &ReservationController{
		reservationUseCase: reservationUseCase,
	}
}

// CreateReservationRequest representa os dados para criar uma reserva
type CreateReservationRequest struct {
	KeyID string `json:"key_id" binding:"required"`
	DueAt string `json:"due_at" binding:"required"` // ISO 8601 format
}

// ExtendReservationRequest representa os dados para estender uma reserva
type ExtendReservationRequest struct {
	NewDueAt string `json:"new_due_at" binding:"required"` // ISO 8601 format
}

// CreateReservation godoc
// @Summary Cria nova reserva
// @Description Cria uma nova reserva de chave
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reservation body CreateReservationRequest true "Dados da reserva"
// @Success 201 {object} response.APIResponse{data=entity.Reservation} "Reserva criada com sucesso"
// @Failure 400 {object} response.APIError "Erro de validação"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 403 {object} response.APIError "Usuário bloqueado"
// @Failure 409 {object} response.APIError "Conflito ao criar reserva"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /reservations [post]
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "Corpo da requisição inválido"))
		return
	}

	keyID, err := primitive.ObjectIDFromHex(req.KeyID)
	if err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "ID da chave inválido"))
		return
	}

	dueAt, err := time.Parse(time.RFC3339, req.DueAt)
	if err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "Data de vencimento inválida"))
		return
	}

	// Obtém o ID do usuário do middleware de autenticação
	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(response.NewError(http.StatusUnauthorized, "Usuário não autenticado"))
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.Error(response.NewError(http.StatusInternalServerError, "ID do usuário inválido"))
		return
	}

	reservation := entity.NewReservation(keyID, userID, dueAt)

	if err := c.reservationUseCase.CreateReservation(ctx, reservation); err != nil {
		switch err {
		case entity.ErrUserBlocked:
			ctx.Error(response.NewError(http.StatusForbidden, "Usuário bloqueado"))
		case entity.ErrKeyInactive:
			ctx.Error(response.NewError(http.StatusConflict, "Chave inativa"))
		case entity.ErrKeyReserved, entity.ErrReservationAlreadyExists:
			ctx.Error(response.NewError(http.StatusConflict, err.Error()))
		default:
			ctx.Error(response.NewError(http.StatusInternalServerError, "Erro ao criar reserva"))
		}
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(reservation, nil))
}

// GetUserReservations godoc
// @Summary Lista reservas do usuário
// @Description Retorna todas as reservas do usuário logado
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]entity.Reservation} "Lista de reservas do usuário"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /reservations [get]
func (c *ReservationController) GetUserReservations(ctx *gin.Context) {
	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(response.NewError(http.StatusUnauthorized, "Usuário não autenticado"))
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.Error(response.NewError(http.StatusInternalServerError, "ID do usuário inválido"))
		return
	}

	reservations, err := c.reservationUseCase.GetUserReservations(ctx, userID)
	if err != nil {
		ctx.Error(response.NewError(http.StatusInternalServerError, "Erro ao buscar reservas"))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(reservations, nil))
}

// GetAllReservations godoc
// @Summary Lista todas as reservas (admin only)
// @Description Retorna uma lista de todas as reservas cadastradas (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]entity.Reservation} "Lista de todas as reservas"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 403 {object} response.APIError "Acesso restrito"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /admin/reservations [get]
func (c *ReservationController) GetAllReservations(ctx *gin.Context) {
	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.Error(response.NewError(http.StatusForbidden, "Acesso restrito a administradores"))
		return
	}

	reservations, err := c.reservationUseCase.GetAllReservations(ctx, userRole.(entity.UserRole))
	if err != nil {
		ctx.Error(response.NewError(http.StatusInternalServerError, "Erro ao buscar reservas"))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(reservations, nil))
}

// ReturnKey godoc
// @Summary Registra devolução de chave
// @Description Marca uma reserva como devolvida
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da reserva"
// @Success 200 {object} response.APIResponse "Chave devolvida com sucesso"
// @Failure 400 {object} response.APIError "ID inválido"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 404 {object} response.APIError "Reserva não encontrada"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /reservations/{id}/return [put]
func (c *ReservationController) ReturnKey(ctx *gin.Context) {
	idParam := ctx.Param("id")
	reservationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "ID da reserva inválido"))
		return
	}

	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(response.NewError(http.StatusUnauthorized, "Usuário não autenticado"))
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.Error(response.NewError(http.StatusInternalServerError, "ID do usuário inválido"))
		return
	}

	if err := c.reservationUseCase.ReturnKey(ctx, reservationID, userID); err != nil {
		if err == entity.ErrReservationNotFound {
			ctx.Error(response.NewError(http.StatusNotFound, "Reserva não encontrada"))
		} else {
			ctx.Error(response.NewError(http.StatusInternalServerError, "Erro ao devolver chave"))
		}
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse("Chave devolvida com sucesso", nil))
}

// GetReservationHistory godoc
// @Summary Histórico de reservas do usuário
// @Description Retorna o histórico de reservas do usuário logado
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]entity.Reservation} "Histórico de reservas"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /reservations/history [get]
func (c *ReservationController) GetReservationHistory(ctx *gin.Context) {
	c.GetUserReservations(ctx)
}

// ExtendReservation godoc
// @Summary Estende prazo de reserva
// @Description Estende o prazo de uma reserva (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da reserva"
// @Param request body ExtendReservationRequest true "Nova data de vencimento"
// @Success 200 {object} response.APIResponse "Reserva estendida com sucesso"
// @Failure 400 {object} response.APIError "Dados inválidos"
// @Failure 401 {object} response.APIError "Não autenticado"
// @Failure 403 {object} response.APIError "Acesso restrito"
// @Failure 404 {object} response.APIError "Reserva não encontrada"
// @Failure 500 {object} response.APIError "Erro interno"
// @Router /admin/reservations/{id}/extend [put]
func (c *ReservationController) ExtendReservation(ctx *gin.Context) {
	idParam := ctx.Param("id")
	reservationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "ID da reserva inválido"))
		return
	}

	var req ExtendReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "Corpo da requisição inválido"))
		return
	}

	newDueAt, err := time.Parse(time.RFC3339, req.NewDueAt)
	if err != nil {
		ctx.Error(response.NewError(http.StatusBadRequest, "Data de vencimento inválida"))
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.Error(response.NewError(http.StatusForbidden, "Acesso restrito a administradores"))
		return
	}

	if err := c.reservationUseCase.ExtendReservation(ctx, reservationID, newDueAt, userRole.(entity.UserRole)); err != nil {
		switch err {
		case entity.ErrReservationNotFound:
			ctx.Error(response.NewError(http.StatusNotFound, "Reserva não encontrada"))
		case entity.ErrCannotExtendReservation:
			ctx.Error(response.NewError(http.StatusBadRequest, err.Error()))
		default:
			ctx.Error(response.NewError(http.StatusInternalServerError, "Erro ao estender reserva"))
		}
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse("Reserva estendida com sucesso", nil))
}
