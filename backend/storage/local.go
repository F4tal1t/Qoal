package storage

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LocalStorage manages local file storage for uploads and processed files
type LocalStorage struct {
	UploadDir    string
	ProcessedDir string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(uploadDir, processedDir string) *LocalStorage {
	// Create directories if they don't exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(processedDir, 0755)

	return &LocalStorage{
		UploadDir:    uploadDir,
		ProcessedDir: processedDir,
	}
}

// SaveFile saves an uploaded file and returns the saved path
func (ls *LocalStorage) SaveFile(file io.Reader, filename string, fileSize int64) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	cleanBaseName := strings.ReplaceAll(baseName, " ", "_")
	cleanBaseName = strings.ReplaceAll(cleanBaseName, "..", "_")

	uniqueID := uuid.New().String()
	savedFilename := fmt.Sprintf("%s_%s%s", cleanBaseName, uniqueID[:8], ext)

	// Create subdirectory based on date for better organization
	dateDir := time.Now().Format("2006/01/02")
	uploadPath := filepath.Join(ls.UploadDir, dateDir)

	// Ensure directory exists
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	fullPath := filepath.Join(uploadPath, savedFilename)

	// Create file
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Copy file content
	written, err := io.Copy(out, file)
	if err != nil {
		os.Remove(fullPath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Verify file size
	if written != fileSize {
		os.Remove(fullPath) // Clean up on error
		return "", fmt.Errorf("file size mismatch: expected %d, got %d", fileSize, written)
	}

	return fullPath, nil
}

// GetFile retrieves a file by path
func (ls *LocalStorage) GetFile(filePath string) (*os.File, error) {
	// Security check: ensure file is within our directories
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %v", err)
	}

	uploadAbs, _ := filepath.Abs(ls.UploadDir)
	processedAbs, _ := filepath.Abs(ls.ProcessedDir)

	if !strings.HasPrefix(absPath, uploadAbs) && !strings.HasPrefix(absPath, processedAbs) {
		return nil, fmt.Errorf("access denied: file outside allowed directories")
	}

	return os.Open(filePath)
}

// DeleteFile removes a file
func (ls *LocalStorage) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// GetOutputPath generates a path for processed files
func (ls *LocalStorage) GetOutputPath(jobID string, targetFormat string) string {
	dateDir := time.Now().Format("2006/01/02")
	outputPath := filepath.Join(ls.ProcessedDir, dateDir)
	os.MkdirAll(outputPath, 0755)

	return filepath.Join(outputPath, fmt.Sprintf("%s.%s", jobID, targetFormat))
}

// ValidateFileType checks if the file type is supported
func ValidateFileType(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	// Supported formats by category
	supportedFormats := map[string][]string{
		"image":    {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".svg"},
		"video":    {".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm", ".m4v"},
		"audio":    {".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma"},
		"document": {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf"},
		"archive":  {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz"},
	}

	for category, formats := range supportedFormats {
		for _, format := range formats {
			if ext == format {
				return category, nil
			}
		}
	}

	return "", fmt.Errorf("unsupported file format: %s", ext)
}

// GetMimeType returns the MIME type for a file extension
func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Fallback for common types
		fallbacks := map[string]string{
			".jpg":  "image/jpeg",
			".jpeg": "image/jpeg",
			".png":  "image/png",
			".gif":  "image/gif",
			".pdf":  "application/pdf",
			".doc":  "application/msword",
			".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			".xls":  "application/vnd.ms-excel",
			".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			".mp3":  "audio/mpeg",
			".mp4":  "video/mp4",
			".zip":  "application/zip",
		}
		if mimeType, ok := fallbacks[ext]; ok {
			return mimeType
		}
		return "application/octet-stream"
	}
	return mimeType
}
