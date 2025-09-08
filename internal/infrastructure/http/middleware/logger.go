package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware logs detailed information about each request.
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

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("| %s | %3d | %13v | %15s | %-7s %s %s\n",
			requestID,
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMsg,
		)
	}
}
