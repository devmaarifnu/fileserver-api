package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/config"
)

// CORSMiddleware configures CORS settings
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	if !cfg.CORS.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           86400, // 24 hours
	}

	// If AllowOrigins contains "*", use AllowAllOrigins
	for _, origin := range cfg.CORS.AllowedOrigins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
			corsConfig.AllowOrigins = nil
			break
		}
	}

	return cors.New(corsConfig)
}
