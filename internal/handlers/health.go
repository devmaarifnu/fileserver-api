package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/services"
)

// HealthHandler handles health check
type HealthHandler struct {
	config         *config.Config
	storageService *services.StorageService
	startTime      time.Time
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(cfg *config.Config, storage *services.StorageService) *HealthHandler {
	return &HealthHandler{
		config:         cfg,
		storageService: storage,
		startTime:      time.Now(),
	}
}

// Handle processes health check request
func (h *HealthHandler) Handle(c *gin.Context) {
	// Calculate uptime
	uptime := time.Since(h.startTime)

	// Get storage info
	storageInfo, err := h.storageService.GetStorageInfo()
	if err != nil {
		storageInfo = map[string]interface{}{
			"error": "Failed to get storage info",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": h.config.App.Version,
		"uptime":  uptime.String(),
		"storage": storageInfo,
	})
}
