package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dessima/gerenciador-chaves-api/entity"
	"github.com/dessima/gerenciador-chaves-api/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// KeyController manipula requisições relacionadas às chaves
type KeyController struct {
	keyUseCase *usecase.KeyUseCase
}

// NewKeyController cria uma nova instância do KeyController
func NewKeyController(keyUseCase *usecase.KeyUseCase) *KeyController {
	return &KeyController{
		keyUseCase: keyUseCase,
	}
}

// TODO: Implementar todos os handlers:

// GetAllKeys godoc
// @Summary Lista todas as chaves
// @Description Retorna uma lista de todas as chaves cadastradas
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Key
// @Failure 401 {object} APIError
// @Failure 500 {object} APIError
// @Router /keys [get]
func (c *KeyController) GetAllKeys(ctx *gin.Context) {
	keys, err := c.keyUseCase.GetAllKeys(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to retrieve keys", Details: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, keys)
}

// GetKeyByID godoc
// @Summary Busca chave por ID
// @Description Retorna uma chave específica por ID
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Key ID"
// @Success 200 {object} entity.Key
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /keys/{id} [get]
func (c *KeyController) GetKeyByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Key ID", Details: err.Error()})
		return
	}

	key, err := c.keyUseCase.GetKeyByID(ctx, id)
	if err != nil {
		if err == entity.ErrKeyNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "Key not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to retrieve key", Details: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, key)
}

// CreateKey godoc
// @Summary Cria nova chave
// @Description Cria uma nova chave (apenas administradores)
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key body entity.Key true "Key data"
// @Success 201 {object} entity.Key
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 500 {object} APIError
// @Router /keys [post]
func (c *KeyController) CreateKey(ctx *gin.Context) {
	var key entity.Key
	if err := ctx.ShouldBindJSON(&key); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}

	// Assuming userRole is set by a middleware
	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.keyUseCase.CreateKey(ctx, &key, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrKeyAlreadyExists {
			ctx.JSON(http.StatusConflict, APIError{Code: http.StatusConflict, Message: "Key with this name already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to create key", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, key)
}

// UpdateKey godoc
// @Summary Atualiza chave existente
// @Description Atualiza dados de uma chave (apenas administradores)
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Key ID"
// @Param key body entity.Key true "Key data"
// @Success 200 {object} entity.Key
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /keys/{id} [put]
func (c *KeyController) UpdateKey(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Key ID", Details: err.Error()})
		return
	}

	var key entity.Key
	if err := ctx.ShouldBindJSON(&key); err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid request body", Details: err.Error()})
		return
	}
	key.ID = id // Ensure the ID from the path is used

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.keyUseCase.UpdateKey(ctx, &key, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrKeyNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "Key not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to update key", Details: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, key)
}

// DeleteKey godoc
// @Summary Remove chave
// @Description Remove uma chave do sistema (apenas administradores)
// @Tags keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Key ID"
// @Success 204
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 403 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /keys/{id} [delete]
func (c *KeyController) DeleteKey(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIError{Code: http.StatusBadRequest, Message: "Invalid Key ID", Details: err.Error()})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
		ctx.JSON(http.StatusForbidden, APIError{Code: http.StatusForbidden, Message: "Admin access required"})
		return
	}

	if err := c.keyUseCase.DeleteKey(ctx, id, userRole.(entity.UserRole)); err != nil {
		if err == entity.ErrKeyNotFound {
			ctx.JSON(http.StatusNotFound, APIError{Code: http.StatusNotFound, Message: "Key not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIError{Code: http.StatusInternalServerError, Message: "Failed to delete key", Details: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
