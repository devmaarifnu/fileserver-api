package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/middleware"
	"github.com/maarifnu/cdn-fileserver/internal/services"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// UploadHandler handles file upload
type UploadHandler struct {
	fileService *services.FileService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(fs *services.FileService) *UploadHandler {
	return &UploadHandler{
		fileService: fs,
	}
}

// Handle processes file upload request
func (h *UploadHandler) Handle(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		logger.WithField("error", err).Warn("No file uploaded")
		utils.ValidationErrorResponse(c, "Validation error", map[string]string{
			"file": "File is required",
		})
		return
	}

	// Get tag
	tag := c.PostForm("tag")
	if tag == "" {
		utils.ValidationErrorResponse(c, "Validation error", map[string]string{
			"tag": "Tag is required",
		})
		return
	}

	// Get public flag (default: false)
	publicStr := c.DefaultPostForm("public", "false")
	public, err := strconv.ParseBool(publicStr)
	if err != nil {
		public = false
	}

	// Get token info from context
	token := middleware.GetTokenFromContext(c)
	tokenName := "Unknown"
	if token != nil {
		tokenName = token.Name
	}

	// Create upload request
	uploadReq := &services.UploadRequest{
		File:       file,
		Tag:        tag,
		Public:     public,
		UploadedBy: tokenName,
	}

	// Upload file
	response, err := h.fileService.Upload(uploadReq)
	if err != nil {
		logger.WithField("error", err).Error("File upload failed")

		// Check if it's a validation error
		if err.Error() == "file is empty" ||
			err.Error()[:8] == "invalid " ||
			err.Error()[:5] == "file " {
			utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", err.Error())
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to upload file")
		return
	}

	logger.WithField("file_id", response.FileID).Info("File uploaded successfully")

	utils.SuccessResponse(c, http.StatusOK, "File uploaded successfully", response)
}
