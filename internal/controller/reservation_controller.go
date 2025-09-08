package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portaria-keys/internal/entity"
	"github.com/portaria-keys/internal/usecase"
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
	KeyID string `json:"key_id" validate:"required"`
	DueAt string `json:"due_at" validate:"required"` // ISO 8601 format
}

// CreateReservation godoc
// @Summary Cria nova reserva
// @Description Cria uma nova reserva de chave
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reservation body CreateReservationRequest true "Reservation data"
// @Success 201 {object} entity.Reservation
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 409 {object} APIError
// @Failure 500 {object} APIError
// @Router /reservations [post]
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}

	keyID, err := primitive.ObjectIDFromHex(req.KeyID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Key ID", Details: err.Error()})
		return
	}

	dueAt, err := time.Parse(time.RFC3339, req.DueAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid DueAt format, use RFC3339", Details: err.Error()})
		return
	}

	// Assuming userID is set by a middleware
	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIError{Code: http.StatusUnauthorized, Message: "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Invalid User ID in context", Details: err.Error()})
		return
	}

	reservation := entity.NewReservation(keyID, userID, dueAt)

	if err := c.reservationUseCase.CreateReservation(ctx, reservation); err != nil {
		if err == entity.ErrUserBlocked || err == entity.ErrKeyInactive || err == entity.ErrKeyReserved || err == entity.ErrReservationAlreadyExists {
			ctx.JSON(http.StatusConflict, APIError{Code: http.StatusConflict, Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to create reservation", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reservation)
}

// GetUserReservations godoc
// @Summary Lista reservas do usuário
// @Description Retorna todas as reservas do usuário logado
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Reservation
// @Failure 401 {object} APIError
// @Failure 500 {object} APIError
// @Router /reservations [get]
func (c *ReservationController) GetUserReservations(ctx *gin.Context) {
	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIError{Code: http.StatusUnauthorized, Message: "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Invalid User ID in context", Details: err.Error()})
		return
	}

	reservations, err := c.reservationUseCase.GetUserReservations(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to retrieve user reservations", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// ReturnKey godoc
// @Summary Registra devolução de chave
// @Description Marca uma reserva como devolvida
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Success 200 {object} entity.Reservation
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /reservations/{id}/return [put]
func (c *ReservationController) ReturnKey(ctx *gin.Context) {
	idParam := ctx.Param("id")
	reservationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Reservation ID", Details: err.Error()})
		return
	}

	userIDHex, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, APIError{Code: http.StatusUnauthorized, Message: "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Invalid User ID in context", Details: err.Error()})
		return
	}

	if err := c.reservationUseCase.ReturnKey(ctx, reservationID, userID); err != nil {
		if err == entity.ErrReservationNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "Reservation not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to return key", Details: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// GetReservationHistory godoc
// @Summary Histórico de reservas do usuário
// @Description Retorna o histórico de reservas do usuário logado
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Reservation
// @Failure 401 {object} APIError
// @Failure 500 {object} APIError
// @Router /reservations/history [get]
func (c *ReservationController) GetReservationHistory(ctx *gin.Context) {
	// This is essentially the same as GetUserReservations, as it returns all reservations for the logged-in user.
	// If a distinction is needed (e.g., only past reservations), additional logic would be required.
	c.GetUserReservations(ctx)
}

// GetAllReservations godoc
// @Summary Lista todas as reservas (admin only)
// @Description Retorna uma lista de todas as reservas cadastradas (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Reservation
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 500 {object} APIError
// @Router /admin/reservations [get]
func (c *ReservationController) GetAllReservations(ctx *gin.Context) {
	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	reservations, err := c.reservationUseCase.GetAllReservations(ctx, userRole.(entity.UserRole))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to retrieve all reservations", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// ExtendReservation godoc
// @Summary Estende prazo de reserva
// @Description Estende o prazo de uma reserva (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Param new_due_at body string true "New due date in RFC3339 format"
// @Success 200 {object} entity.Reservation
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /admin/reservations/{id}/extend [put]
func (c *ReservationController) ExtendReservation(ctx *gin.Context) {
	idParam := ctx.Param("id")
	reservationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Reservation ID", Details: err.Error()})
		return
	}

	var req struct {
		NewDueAt string `json:"new_due_at" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}

	newDueAt, err := time.Parse(time.RFC3339, req.NewDueAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid NewDueAt format, use RFC3339", Details: err.Error()})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.reservationUseCase.ExtendReservation(ctx, reservationID, newDueAt, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrReservationNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "Reservation not found"})
			return
		}
		if err == entity.ErrCannotExtendReservation {
			ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to extend reservation", Details: err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
