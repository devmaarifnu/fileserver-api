package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()

		// Get request details
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log incoming request
		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     method,
			"path":       path,
			"ip":         clientIP,
			"user_agent": userAgent,
		}).Info("Incoming request")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response details
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		// Determine log level based on status code
		logEntry := logger.WithFields(logrus.Fields{
			"request_id":    requestID,
			"method":        method,
			"path":          path,
			"status":        statusCode,
			"latency":       latency.String(),
			"latency_ms":    latency.Milliseconds(),
			"response_size": responseSize,
			"ip":            clientIP,
		})

		// Log response
		if statusCode >= 500 {
			logEntry.Error("Request completed with server error")
		} else if statusCode >= 400 {
			logEntry.Warn("Request completed with client error")
		} else {
			logEntry.Info("Request completed successfully")
		}
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}
