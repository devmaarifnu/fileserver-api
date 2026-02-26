package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// tagRegex allows alphanumeric, dash, and underscore
	tagRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidateTag validates a tag name
func ValidateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag is required")
	}

	if len(tag) > 50 {
		return fmt.Errorf("tag is too long (max 50 characters)")
	}

	if !tagRegex.MatchString(tag) {
		return fmt.Errorf("tag must contain only alphanumeric characters, dashes, and underscores")
	}

	return nil
}

// ValidateFileExtension checks if file extension is allowed
func ValidateFileExtension(filename string, allowedExtensions []string) error {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	ext = strings.ToLower(ext)

	if ext == "" {
		return fmt.Errorf("file has no extension")
	}

	for _, allowed := range allowedExtensions {
		if strings.ToLower(allowed) == ext {
			return nil
		}
	}

	return fmt.Errorf("file extension '.%s' is not allowed", ext)
}

// ValidateFileSize checks if file size is within limit
func ValidateFileSize(size int64, maxSize int64) error {
	if size <= 0 {
		return fmt.Errorf("file is empty")
	}

	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum limit (%d MB)", maxSize/1024/1024)
	}

	return nil
}

// SanitizeFilename removes potentially dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators
	filename = filepath.Base(filename)

	// Remove any null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")

	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove any characters that are not alphanumeric, dash, underscore, or dot
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = reg.ReplaceAllString(filename, "")

	// Remove leading dots (hidden files)
	filename = strings.TrimPrefix(filename, ".")

	// Limit length
	if len(filename) > 200 {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		name = name[:200-len(ext)]
		filename = name + ext
	}

	return filename
}

// IsValidFilename checks if filename is valid
func IsValidFilename(filename string) bool {
	if filename == "" {
		return false
	}

	if strings.Contains(filename, "..") {
		return false
	}

	if strings.ContainsAny(filename, "/\\") {
		return false
	}

	return true
}
