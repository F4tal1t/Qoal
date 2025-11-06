package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Max upload size: 100MB
	MaxFileSize = 100 * 1024 * 1024
	// Upload directory
	UploadDir = "uploads"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler() *UploadHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic("Failed to create upload directory: " + err.Error())
	}

	return &UploadHandler{
		uploadDir: UploadDir,
	}
}

type UploadResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	FilePath  string `json:"file_path,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	FileSize  int64  `json:"file_size,omitempty"`
	UploadID  string `json:"upload_id,omitempty"`
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	// 🎓 **Go Tutorial: Multipart Form Handling**
	// This is how we handle file uploads in Go
	
	// Limit upload size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "File too large or invalid form data",
		})
		return
	}

	// Get file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "No file provided",
		})
		return
	}
	defer file.Close()

	// Validate file type
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedFileType(fileExt) {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "File type not supported: " + fileExt,
		})
		return
	}

	// Generate unique filename
	uploadID := uuid.New().String()
	newFilename := fmt.Sprintf("%s%s", uploadID, fileExt)
	filePath := filepath.Join(h.uploadDir, newFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Failed to create file on server",
		})
		return
	}
	defer dst.Close()

	// Copy file contents
	written, err := io.Copy(dst, file)
	if err != nil {
		// Clean up on error
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	// 🎓 **Go Tutorial: Defer Statements**
	// defer ensures cleanup happens even if function panics
	// defer dst.Close() was called earlier

	c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		Message:  "File uploaded successfully",
		FilePath: filePath,
		FileName: header.Filename,
		FileSize: written,
		UploadID: uploadID,
	})
}

// 🎓 **Go Tutorial: Slice Literals and Functions**
func isAllowedFileType(ext string) bool {
	// Supported file types for our processor
	allowedTypes := []string{
		// Images
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".heic",
		// Documents
		".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt",
		// Audio
		".mp3", ".wav", ".flac", ".aac", ".m4a", ".ogg",
		// Video
		".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv", ".wmv",
		// Archives
		".zip", ".rar", ".7z", ".tar.gz",
	}

	// 🎓 **Go Tutorial: Range Loops**
	for _, allowed := range allowedTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (h *UploadHandler) HandleListUploads(c *gin.Context) {
	// 🎓 **Go Tutorial: Reading Directories**
	files, err := os.ReadDir(h.uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read upload directory",
		})
		return
	}

	var uploads []map[string]interface{}
	
	// 🎓 **Go Tutorial: Working with File Info**
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err == nil {
				uploads = append(uploads, map[string]interface{}{
					"name":      info.Name(),
					"size":      info.Size(),
					"modified":  info.ModTime().Format(time.RFC3339),
					"extension": filepath.Ext(info.Name()),
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"uploads": uploads,
		"count":   len(uploads),
	})
}

func (h *UploadHandler) HandleDeleteUpload(c *gin.Context) {
	filename := c.Param("filename")
	
	// Security: Prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filename",
		})
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)
	
	// 🎓 **Go Tutorial: File Operations**
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete file",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
		"file":    filename,
	})
}