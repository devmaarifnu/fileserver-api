package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/services"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// ListHandler handles file listing
type ListHandler struct {
	fileService *services.FileService
	config      *config.Config
}

// NewListHandler creates a new list handler
func NewListHandler(fs *services.FileService, cfg *config.Config) *ListHandler {
	return &ListHandler{
		fileService: fs,
		config:      cfg,
	}
}

// Handle processes file list request
func (h *ListHandler) Handle(c *gin.Context) {
	// Parse query parameters
	tag := c.Query("tag")
	search := c.Query("search")

	// Parse public filter
	var publicFilter *bool
	if publicStr := c.Query("public"); publicStr != "" {
		if public, err := strconv.ParseBool(publicStr); err == nil {
			publicFilter = &public
		}
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Parse sort order
	sortOrder := c.DefaultQuery("sort", "desc")
	sortDesc := sortOrder == "desc"

	// Create list request
	listReq := &services.ListRequest{
		Tag:      tag,
		Public:   publicFilter,
		Search:   search,
		Page:     page,
		Limit:    limit,
		SortDesc: sortDesc,
	}

	// Get files
	files, totalItems, err := h.fileService.List(listReq)
	if err != nil {
		logger.WithField("error", err).Error("Failed to list files")
		utils.InternalServerErrorResponse(c, "Failed to retrieve files")
		return
	}

	// Build file responses with URLs
	baseURL := h.config.GetBaseURL()
	fileResponses := make([]map[string]interface{}, 0, len(files))

	for _, file := range files {
		fileResponses = append(fileResponses, map[string]interface{}{
			"file_id":       file.FileID,
			"original_name": file.OriginalName,
			"tag":           file.Tag,
			"url":           baseURL + "/" + file.Tag + "/" + file.FileID,
			"size":          file.Size,
			"content_type":  file.ContentType,
			"public":        file.Public,
			"uploaded_at":   file.UploadedAt,
			"uploaded_by":   file.UploadedBy,
		})
	}

	// Calculate pagination metadata
	totalPages := (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	pagination := &utils.PaginationMeta{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}

	utils.PaginationResponse(c, "Files retrieved successfully", fileResponses, pagination)
}
