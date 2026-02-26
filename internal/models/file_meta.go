package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileMeta represents file metadata
type FileMeta struct {
	FileID       string    `json:"file_id"`
	OriginalName string    `json:"original_name"`
	Tag          string    `json:"tag"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	Public       bool      `json:"public"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UploadedBy   string    `json:"uploaded_by"`
}

// Save saves metadata to a JSON file
func (fm *FileMeta) Save(storagePath string) error {
	metaPath := fm.GetMetaPath(storagePath)

	// Ensure directory exists
	metaDir := filepath.Dir(metaPath)
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(fm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write to file
	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// Load loads metadata from a JSON file
func (fm *FileMeta) Load(storagePath string) error {
	metaPath := fm.GetMetaPath(storagePath)

	// Read file
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, fm); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}

// Delete deletes the metadata file
func (fm *FileMeta) Delete(storagePath string) error {
	metaPath := fm.GetMetaPath(storagePath)

	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata file: %w", err)
	}

	return nil
}

// Exists checks if metadata file exists
func (fm *FileMeta) Exists(storagePath string) bool {
	metaPath := fm.GetMetaPath(storagePath)
	_, err := os.Stat(metaPath)
	return err == nil
}

// GetMetaPath returns the path to the metadata file
func (fm *FileMeta) GetMetaPath(storagePath string) string {
	return filepath.Join(storagePath, fm.Tag, fm.FileID+".meta.json")
}

// GetFilePath returns the path to the actual file
func (fm *FileMeta) GetFilePath(storagePath string) string {
	return filepath.Join(storagePath, fm.Tag, fm.FileID)
}

// LoadFromFile loads metadata from a specific file path
func LoadFromFile(metaPath string) (*FileMeta, error) {
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var meta FileMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &meta, nil
}
