package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/portaria-keys/internal/controller"
	"github.com/portaria-keys/internal/entity"
)

// TODO: Mover para configuração
const jwtSecret = "supersecretkey"

// AuthMiddleware valida o token JWT e extrai informações do usuário
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, controller.APIError{Code: http.StatusUnauthorized, Message: "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, controller.APIError{Code: http.StatusUnauthorized, Message: "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, controller.APIError{Code: http.StatusUnauthorized, Message: "Invalid token", Details: err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusUnauthorized, controller.APIError{Code: http.StatusUnauthorized, Message: "Token expired"})
				c.Abort()
				return
			}

			c.Set("userID", claims["user_id"])
			c.Set("userRole", entity.UserRole(claims["role"].(string)))
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, controller.APIError{Code: http.StatusUnauthorized, Message: "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// AdminMiddleware verifica se o usuário autenticado é um administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || userRole.(entity.UserRole) != entity.UserRoleAdmin {
			c.JSON(http.StatusForbidden, controller.APIError{Code: http.StatusForbidden, Message: "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}