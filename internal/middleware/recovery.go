package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()

				// Get request ID
				requestID := GetRequestID(c)

				// Log the panic
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      err,
					"stack":      string(stack),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"ip":         c.ClientIP(),
				}).Error("Panic recovered")

				// Send error response
				utils.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Internal Server Error",
					fmt.Sprintf("An unexpected error occurred: %v", err),
				)

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}
