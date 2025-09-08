package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/portaria-keys/internal/controller"
)

// RateLimitMiddleware implementa um rate limiting simples usando um token bucket.
// Limita a 10 requisições por segundo, com um burst de 5.
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second/10), 5) // 10 req/s, burst 5

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, controller.APIError{Code: http.StatusTooManyRequests, Message: "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}