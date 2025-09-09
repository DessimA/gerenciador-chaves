package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecurityEvent representa um evento de seguran�a
type SecurityEvent struct {
	Type      string    `json:"type"`
	UserID    string    `json:"user_id,omitempty"`
	IP        string    `json:"ip"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// LoggerMiddleware logs detailed information about each request with enhanced security event tracking.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set("requestID", requestID)

		// Process request
		c.Next()

		// Log after request is processed
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		userID, _ := c.Get("user_id")

		if raw != "" {
			path = path + "?" + raw
		}

		// Log basic request info
		log.Printf("| %s | %3d | %13v | %15s | %-7s %s %s\n",
			requestID,
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMsg,
		)

		// Log security events for specific conditions
		if isSecurityEvent(statusCode, path, method) {
			event := SecurityEvent{
				Type:      determineEventType(statusCode, path, method),
				UserID:    formatUserID(userID),
				IP:        clientIP,
				Method:    method,
				Path:      path,
				Status:    statusCode,
				Error:     errorMsg,
				Timestamp: end,
			}
			logSecurityEvent(event)
		}
	}
}

// isSecurityEvent determina se um evento deve ser registrado como evento de seguran�a
func isSecurityEvent(status int, path, method string) bool {
	// Falhas de autentica��o/autoriza��o
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return true
	}

	// Tentativas de login/registro
	if path == "/api/v1/auth/login" || path == "/api/v1/auth/register" {
		return true
	}

	// A��es administrativas
	if strings.HasPrefix(path, "/api/v1/admin/") {
		return true
	}

	return false
}

// determineEventType retorna o tipo do evento de seguran�a
func determineEventType(status int, path, method string) string {
	if status == http.StatusUnauthorized {
		return "UNAUTHORIZED_ACCESS"
	}
	if status == http.StatusForbidden {
		return "FORBIDDEN_ACCESS"
	}
	if path == "/api/v1/auth/login" {
		return "LOGIN_ATTEMPT"
	}
	if path == "/api/v1/auth/register" {
		return "REGISTER_ATTEMPT"
	}
	if strings.HasPrefix(path, "/api/v1/admin/") {
		return "ADMIN_ACTION"
	}
	return "SECURITY_EVENT"
}

// formatUserID formata o ID do usu�rio para logging
func formatUserID(userID interface{}) string {
	if userID == nil {
		return ""
	}
	return fmt.Sprint(userID)
}

// logSecurityEvent registra um evento de seguran�a
func logSecurityEvent(event SecurityEvent) {
	eventJSON, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Printf("Erro ao serializar evento de seguran�a: %v", err)
		return
	}
	log.Printf("Evento de Seguran�a:\n%s", string(eventJSON))
}
