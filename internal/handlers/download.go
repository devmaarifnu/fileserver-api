package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/middleware"
	"github.com/maarifnu/cdn-fileserver/internal/services"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// DownloadHandler handles file download/view
type DownloadHandler struct {
	fileService *services.FileService
}

// NewDownloadHandler creates a new download handler
func NewDownloadHandler(fs *services.FileService) *DownloadHandler {
	return &DownloadHandler{
		fileService: fs,
	}
}

// Handle processes file download/view request
func (h *DownloadHandler) Handle(c *gin.Context) {
	tag := c.Param("tag")
	filename := c.Param("filename")

	// Validate parameters
	if tag == "" || filename == "" {
		utils.NotFoundResponse(c, "File not found")
		return
	}

	// Get file metadata
	meta, filePath, err := h.fileService.Download(tag, filename)
	if err != nil {
		logger.WithField("error", err).Warn("File not found")
		utils.NotFoundResponse(c, "File not found")
		return
	}

	// Check if file is private
	if !meta.Public {
		// Check authentication
		isAuth := middleware.IsAuthenticated(c)
		if !isAuth {
			logger.WithField("file_id", filename).Warn("Unauthorized access to private file")
			utils.ForbiddenResponse(c, "This file is private and requires authentication")
			return
		}
	}

	// Check if download parameter is set
	download := c.Query("download")
	if download == "true" {
		c.Header("Content-Disposition", "attachment; filename=\""+meta.OriginalName+"\"")
	}

	// Set Content-Type header
	c.Header("Content-Type", meta.ContentType)

	// Set cache headers for public files
	if meta.Public {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	}

	// Serve file
	c.File(filePath)

	logger.WithField("file_id", filename).Debug("File served successfully")
}
