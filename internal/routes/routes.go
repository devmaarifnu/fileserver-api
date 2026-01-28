package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/handlers"
	"github.com/maarifnu/cdn-fileserver/internal/middleware"
	"github.com/maarifnu/cdn-fileserver/internal/services"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	storageService *services.StorageService,
	fileService *services.FileService,
) {
	// Create handlers
	uploadHandler := handlers.NewUploadHandler(fileService)
	downloadHandler := handlers.NewDownloadHandler(fileService)
	listHandler := handlers.NewListHandler(fileService, cfg)
	deleteHandler := handlers.NewDeleteHandler(fileService)
	healthHandler := handlers.NewHealthHandler(cfg, storageService)

	// Apply global middleware
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))

	// Security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// Public routes
	router.GET("/health", healthHandler.Handle)

	// File download/view route with optional authentication
	router.GET("/:tag/:filename", middleware.OptionalAuth(cfg), downloadHandler.Handle)

	// API group - requires authentication
	api := router.Group("/api")
	{
		// File management routes
		files := api.Group("/files")
		{
			// List files - requires list permission
			files.GET("", middleware.TokenAuth(cfg, "list"), listHandler.Handle)

			// Delete file - requires delete permission
			files.DELETE("/:tag/:filename", middleware.TokenAuth(cfg, "delete"), deleteHandler.Handle)
		}
	}

	// Upload route - requires upload permission
	router.POST("/upload", middleware.TokenAuth(cfg, "upload"), uploadHandler.Handle)
}
