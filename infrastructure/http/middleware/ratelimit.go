package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/dessima/gerenciador-chaves-api/controller"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/config"
)

// userLimiter mantém um mapa de limitadores por usuário
type userLimiter struct {
	sync.RWMutex
	limiters map[string]*rate.Limiter
	config   *config.Config
}

// newUserLimiter cria um novo gerenciador de rate limiting por usuário
func newUserLimiter(cfg *config.Config) *userLimiter {
	return &userLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   cfg,
	}
}

// getLimiter retorna um limitador para um usuário específico
func (ul *userLimiter) getLimiter(userID string) *rate.Limiter {
	ul.RLock()
	limiter, exists := ul.limiters[userID]
	ul.RUnlock()

	if !exists {
		ul.Lock()
		// Double check para evitar race conditions
		if limiter, exists = ul.limiters[userID]; !exists {
			// 100 requisições por minuto por usuário por padrão
			limiter = rate.NewLimiter(rate.Every(time.Minute/100), 10)
			ul.limiters[userID] = limiter
		}
		ul.Unlock()
	}

	return limiter
}

// RateLimitMiddleware implementa um rate limiting por usuário usando token bucket.
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	ul := newUserLimiter(cfg)
	// Limitador global para IPs não autenticados
	globalLimiter := rate.NewLimiter(rate.Every(time.Second), 3) // 3 req/s para não autenticados

	return func(c *gin.Context) {
		// Tenta obter o ID do usuário do contexto (definido pelo AuthMiddleware)
		userID, exists := c.Get("user_id")

		var limiter *rate.Limiter
		if exists {
			// Usa o limitador específico do usuário se estiver autenticado
			limiter = ul.getLimiter(userID.(string))
		} else {
			// Usa o limitador global para requisições não autenticadas
			limiter = globalLimiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, controller.APIError{
				Code:    http.StatusTooManyRequests,
				Message: "Muitas requisições. Por favor, tente novamente em alguns instantes.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
