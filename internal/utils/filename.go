package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// GenerateUniqueFilename generates a unique filename with UUID
// Format: {sanitized_name}_{uuid}.{ext}
func GenerateUniqueFilename(originalName string) string {
	// Sanitize the original name
	sanitized := SanitizeFilename(originalName)

	// Extract extension
	ext := ExtractExtension(sanitized)

	// Extract name without extension
	nameWithoutExt := strings.TrimSuffix(sanitized, ext)

	// Remove trailing dots
	nameWithoutExt = strings.TrimSuffix(nameWithoutExt, ".")

	// If name is empty, use "file"
	if nameWithoutExt == "" {
		nameWithoutExt = "file"
	}

	// Generate UUID
	id := uuid.New().String()

	// Take only first 8 characters of UUID for shorter filenames
	shortID := strings.Split(id, "-")[0]

	// Combine: name_uuid.ext
	if ext != "" {
		return fmt.Sprintf("%s_%s%s", nameWithoutExt, shortID, ext)
	}

	return fmt.Sprintf("%s_%s", nameWithoutExt, shortID)
}

// ExtractExtension extracts file extension including the dot
func ExtractExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

// SanitizeName cleans the filename (same as SanitizeFilename)
func SanitizeName(filename string) string {
	return SanitizeFilename(filename)
}

// GetNameWithoutExtension returns filename without extension
func GetNameWithoutExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}
