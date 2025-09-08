package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portaria-keys/internal/entity"
	"github.com/portaria-keys/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminController manipula requisições administrativas
type AdminController struct {
	userUseCase *usecase.UserUseCase
}

// NewAdminController cria uma nova instância do AdminController
func NewAdminController(userUseCase *usecase.UserUseCase) *AdminController {
	return &AdminController{
		userUseCase: userUseCase,
	}
}

// BlockUser godoc
// @Summary Bloqueia um usuário (admin only)
// @Description Bloqueia um usuário específico pelo ID (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {string} string "User blocked successfully"
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /admin/users/{id}/block [post]
func (c *AdminController) BlockUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid User ID", Details: err.Error()})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.userUseCase.BlockUser(ctx, userID, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to block user", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// UnblockUser godoc
// @Summary Desbloqueia um usuário (admin only)
// @Description Desbloqueia um usuário específico pelo ID (apenas administradores)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {string} string "User unblocked successfully"
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /admin/users/{id}/unblock [post]
func (c *AdminController) UnblockUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid User ID", Details: err.Error()})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.userUseCase.UnblockUser(ctx, userID, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to unblock user", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}