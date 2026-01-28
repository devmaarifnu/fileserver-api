package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/services"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// DeleteHandler handles file deletion
type DeleteHandler struct {
	fileService *services.FileService
}

// NewDeleteHandler creates a new delete handler
func NewDeleteHandler(fs *services.FileService) *DeleteHandler {
	return &DeleteHandler{
		fileService: fs,
	}
}

// Handle processes file delete request
func (h *DeleteHandler) Handle(c *gin.Context) {
	tag := c.Param("tag")
	filename := c.Param("filename")

	// Validate parameters
	if tag == "" || filename == "" {
		utils.NotFoundResponse(c, "File not found")
		return
	}

	// Delete file
	err := h.fileService.Delete(tag, filename)
	if err != nil {
		if err.Error() == "file not found" {
			utils.NotFoundResponse(c, "File not found")
			return
		}

		logger.WithField("error", err).Error("Failed to delete file")
		utils.InternalServerErrorResponse(c, "Failed to delete file")
		return
	}

	logger.WithField("file_id", filename).Info("File deleted successfully")

	utils.SuccessResponse(c, http.StatusOK, "File deleted successfully", map[string]interface{}{
		"file_id":    filename,
		"tag":        tag,
		"deleted_at": time.Now(),
	})
}
