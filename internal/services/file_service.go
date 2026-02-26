package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sort"
	"time"

	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/models"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// FileService handles file operations
type FileService struct {
	config         *config.Config
	storageService *StorageService
}

// NewFileService creates a new file service
func NewFileService(cfg *config.Config, storage *StorageService) *FileService {
	return &FileService{
		config:         cfg,
		storageService: storage,
	}
}

// UploadRequest represents a file upload request
type UploadRequest struct {
	File         *multipart.FileHeader
	Tag          string
	Public       bool
	UploadedBy   string
}

// UploadResponse represents a file upload response
type UploadResponse struct {
	FileID       string    `json:"file_id"`
	OriginalName string    `json:"original_name"`
	URL          string    `json:"url"`
	Tag          string    `json:"tag"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	Public       bool      `json:"public"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UploadedBy   string    `json:"uploaded_by"`
}

// Upload handles file upload
func (fs *FileService) Upload(req *UploadRequest) (*UploadResponse, error) {
	// Validate tag
	if err := utils.ValidateTag(req.Tag); err != nil {
		return nil, fmt.Errorf("invalid tag: %w", err)
	}

	// Validate file size
	if err := utils.ValidateFileSize(req.File.Size, fs.config.Storage.MaxFileSize); err != nil {
		return nil, err
	}

	// Validate file extension
	if err := utils.ValidateFileExtension(req.File.Filename, fs.config.Storage.AllowedExtensions); err != nil {
		return nil, err
	}

	// Generate unique filename
	fileID := utils.GenerateUniqueFilename(req.File.Filename)

	// Open uploaded file
	src, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Save file to storage
	if err := fs.storageService.SaveFile(req.Tag, fileID, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Get file path for content type detection
	filePath, err := fs.storageService.GetFile(req.Tag, fileID)
	if err != nil {
		// Cleanup on error
		fs.storageService.DeleteFile(req.Tag, fileID)
		return nil, fmt.Errorf("failed to get file path: %w", err)
	}

	// Detect content type
	contentType, err := utils.GetContentType(filePath)
	if err != nil {
		logger.Warnf("Failed to detect content type: %v", err)
		contentType = "application/octet-stream"
	}

	// Create metadata
	meta := &models.FileMeta{
		FileID:       fileID,
		OriginalName: req.File.Filename,
		Tag:          req.Tag,
		Size:         req.File.Size,
		ContentType:  contentType,
		Public:       req.Public,
		UploadedAt:   time.Now(),
		UploadedBy:   req.UploadedBy,
	}

	// Save metadata
	if err := meta.Save(fs.config.Storage.BasePath); err != nil {
		// Cleanup on error
		fs.storageService.DeleteFile(req.Tag, fileID)
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	// Build URL
	baseURL := fs.config.GetBaseURL()
	fileURL := fmt.Sprintf("%s/%s/%s", baseURL, req.Tag, fileID)

	logger.WithField("file_id", fileID).Info("File uploaded successfully")

	return &UploadResponse{
		FileID:       fileID,
		OriginalName: req.File.Filename,
		URL:          fileURL,
		Tag:          req.Tag,
		Size:         req.File.Size,
		ContentType:  contentType,
		Public:       req.Public,
		UploadedAt:   meta.UploadedAt,
		UploadedBy:   req.UploadedBy,
	}, nil
}

// Download retrieves file metadata for download
func (fs *FileService) Download(tag, fileID string) (*models.FileMeta, string, error) {
	// Load metadata
	meta := &models.FileMeta{
		Tag:    tag,
		FileID: fileID,
	}

	if err := meta.Load(fs.config.Storage.BasePath); err != nil {
		return nil, "", fmt.Errorf("file not found")
	}

	// Get file path
	filePath, err := fs.storageService.GetFile(tag, fileID)
	if err != nil {
		return nil, "", fmt.Errorf("file not found")
	}

	return meta, filePath, nil
}

// ListRequest represents a file list request
type ListRequest struct {
	Tag      string
	Public   *bool
	Search   string
	Page     int
	Limit    int
	SortDesc bool
}

// List retrieves a list of files with pagination
func (fs *FileService) List(req *ListRequest) ([]*models.FileMeta, int, error) {
	// Get all files matching filters
	allFiles, err := fs.storageService.ListFiles(req.Tag, req.Public, req.Search)
	if err != nil {
		return nil, 0, err
	}

	// Sort by upload date
	sort.Slice(allFiles, func(i, j int) bool {
		if req.SortDesc {
			return allFiles[i].UploadedAt.After(allFiles[j].UploadedAt)
		}
		return allFiles[i].UploadedAt.Before(allFiles[j].UploadedAt)
	})

	// Calculate pagination
	totalItems := len(allFiles)
	startIndex := (req.Page - 1) * req.Limit
	endIndex := startIndex + req.Limit

	// Handle out of range
	if startIndex >= totalItems {
		return []*models.FileMeta{}, totalItems, nil
	}

	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Return paginated slice
	paginatedFiles := allFiles[startIndex:endIndex]

	return paginatedFiles, totalItems, nil
}

// Delete removes a file and its metadata
func (fs *FileService) Delete(tag, fileID string) error {
	// Load metadata first to check if file exists
	meta := &models.FileMeta{
		Tag:    tag,
		FileID: fileID,
	}

	if err := meta.Load(fs.config.Storage.BasePath); err != nil {
		return fmt.Errorf("file not found")
	}

	// Delete actual file
	if err := fs.storageService.DeleteFile(tag, fileID); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete metadata
	if err := meta.Delete(fs.config.Storage.BasePath); err != nil {
		logger.Warnf("Failed to delete metadata: %v", err)
	}

	logger.WithField("file_id", fileID).Info("File deleted successfully")

	return nil
}

// GetFile returns file reader for streaming
func (fs *FileService) GetFile(tag, fileID string) (io.ReadCloser, error) {
	filePath, err := fs.storageService.GetFile(tag, fileID)
	if err != nil {
		return nil, err
	}

	file, err := openFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Helper function to open file
func openFile(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
