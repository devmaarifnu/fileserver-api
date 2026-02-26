package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/maarifnu/cdn-fileserver/internal/config"
	"github.com/maarifnu/cdn-fileserver/internal/models"
	"github.com/maarifnu/cdn-fileserver/internal/utils"
	"github.com/maarifnu/cdn-fileserver/pkg/logger"
)

// StorageService handles file storage operations
type StorageService struct {
	config *config.Config
}

// NewStorageService creates a new storage service
func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{
		config: cfg,
	}
}

// SaveFile saves a file to storage
func (s *StorageService) SaveFile(tag, fileID string, src io.Reader) error {
	// Create directory if it doesn't exist
	dirPath := filepath.Join(s.config.Storage.BasePath, tag)
	if err := utils.CreateDirectory(dirPath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create destination file
	filePath := filepath.Join(dirPath, fileID)
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		// Clean up on error
		os.Remove(filePath)
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Set file permissions
	if err := os.Chmod(filePath, 0644); err != nil {
		logger.Warnf("Failed to set file permissions: %v", err)
	}

	return nil
}

// GetFile retrieves a file from storage
func (s *StorageService) GetFile(tag, fileID string) (string, error) {
	filePath := filepath.Join(s.config.Storage.BasePath, tag, fileID)

	if !utils.FileExists(filePath) {
		return "", fmt.Errorf("file not found")
	}

	return filePath, nil
}

// DeleteFile deletes a file from storage
func (s *StorageService) DeleteFile(tag, fileID string) error {
	filePath := filepath.Join(s.config.Storage.BasePath, tag, fileID)

	if err := utils.DeleteFile(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListFiles lists all files in storage with optional filters
func (s *StorageService) ListFiles(filterTag string, filterPublic *bool, search string) ([]*models.FileMeta, error) {
	var files []*models.FileMeta

	// Walk through storage directory
	err := filepath.Walk(s.config.Storage.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .meta.json files
		if !strings.HasSuffix(path, ".meta.json") {
			return nil
		}

		// Load metadata
		meta, err := models.LoadFromFile(path)
		if err != nil {
			logger.Warnf("Failed to load metadata from %s: %v", path, err)
			return nil // Skip this file
		}

		// Apply filters
		if filterTag != "" && meta.Tag != filterTag {
			return nil
		}

		if filterPublic != nil && meta.Public != *filterPublic {
			return nil
		}

		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(meta.OriginalName), searchLower) &&
				!strings.Contains(strings.ToLower(meta.FileID), searchLower) {
				return nil
			}
		}

		files = append(files, meta)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// GetStorageInfo returns storage statistics
func (s *StorageService) GetStorageInfo() (map[string]interface{}, error) {
	totalFiles := 0
	var totalSize int64

	// Walk through storage directory
	err := filepath.Walk(s.config.Storage.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip metadata files
		if strings.HasSuffix(path, ".meta.json") {
			return nil
		}

		// Skip .gitkeep files
		if info.Name() == ".gitkeep" {
			return nil
		}

		totalFiles++
		totalSize += info.Size()
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get storage info: %w", err)
	}

	return map[string]interface{}{
		"total_files": totalFiles,
		"total_size":  utils.FormatFileSize(totalSize),
	}, nil
}

// FileExists checks if a file exists
func (s *StorageService) FileExists(tag, fileID string) bool {
	filePath := filepath.Join(s.config.Storage.BasePath, tag, fileID)
	return utils.FileExists(filePath)
}
