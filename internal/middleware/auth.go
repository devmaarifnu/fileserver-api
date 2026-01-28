package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
	"github.com/sirupsen/logrus"
)

// TokenAuth is a middleware for token-based authentication
func TokenAuth(cfg *config.Config, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token := extractTokenFromHeader(c)

		// If not found in header, try query parameter (for file downloads)
		if token == "" {
			token = c.Query("token")
		}

		// If no token provided
		if token == "" {
			logger.WithFields(logrus.Fields{
				"ip":   c.ClientIP(),
				"path": c.Request.URL.Path,
			}).Warn("Missing authentication token")

			utils.UnauthorizedResponse(c, "Invalid or missing token")
			return
		}

		// Find token in config
		tokenConfig := cfg.FindTokenByKey(token)
		if tokenConfig == nil {
			logger.WithFields(logrus.Fields{
				"ip":   c.ClientIP(),
				"path": c.Request.URL.Path,
			}).Warn("Invalid authentication token")

			utils.UnauthorizedResponse(c, "Invalid or missing token")
			return
		}

		// Check if token has required permission
		if requiredPermission != "" && !tokenConfig.HasPermission(requiredPermission) {
			logger.WithFields(logrus.Fields{
				"ip":         c.ClientIP(),
				"path":       c.Request.URL.Path,
				"token_name": tokenConfig.Name,
				"required":   requiredPermission,
			}).Warn("Insufficient permissions")

			utils.ForbiddenResponse(c, "Token does not have "+requiredPermission+" permission")
			return
		}

		// Store token info in context
		c.Set("token", tokenConfig)
		c.Set("token_name", tokenConfig.Name)

		logger.WithFields(logrus.Fields{
			"token_name": tokenConfig.Name,
			"path":       c.Request.URL.Path,
		}).Debug("Authentication successful")

		c.Next()
	}
}

// OptionalAuth is middleware for optional authentication (for public/private file access)
func OptionalAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token := extractTokenFromHeader(c)

		// If not found in header, try query parameter
		if token == "" {
			token = c.Query("token")
		}

		// If token provided, validate it
		if token != "" {
			tokenConfig := cfg.FindTokenByKey(token)
			if tokenConfig != nil {
				// Store token info in context
				c.Set("token", tokenConfig)
				c.Set("token_name", tokenConfig.Name)
				c.Set("authenticated", true)

				logger.WithFields(logrus.Fields{
					"token_name": tokenConfig.Name,
					"path":       c.Request.URL.Path,
				}).Debug("Optional authentication successful")
			}
		} else {
			c.Set("authenticated", false)
		}

		c.Next()
	}
}

// extractTokenFromHeader extracts the token from the Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for "Bearer " prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

// GetTokenFromContext retrieves token config from context
func GetTokenFromContext(c *gin.Context) *config.TokenConfig {
	if token, exists := c.Get("token"); exists {
		if tokenConfig, ok := token.(*config.TokenConfig); ok {
			return tokenConfig
		}
	}
	return nil
}

// IsAuthenticated checks if the request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	if auth, exists := c.Get("authenticated"); exists {
		if authenticated, ok := auth.(bool); ok {
			return authenticated
		}
	}
	return GetTokenFromContext(c) != nil
}
