package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/storage"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100MB
)

type UploadHandler struct {
	db           *gorm.DB
	localStorage *storage.LocalStorage
	jobService   *services.JobService
}

func NewUploadHandler(db *gorm.DB, localStorage *storage.LocalStorage, jobService *services.JobService) *UploadHandler {
	return &UploadHandler{
		db:           db,
		localStorage: localStorage,
		jobService:   jobService,
	}
}

type UploadRequest struct {
	TargetFormat  string `form:"target_format" binding:"required"`
	QualityPreset string `form:"quality_preset"`
}

type UploadResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	JobID        string                 `json:"job_id,omitempty"`
	Status       string                 `json:"status,omitempty"`
	FileInfo     map[string]interface{} `json:"file_info,omitempty"`
	OriginalName string                 `json:"original_filename"`
	FileSize     int64                  `json:"file_size"`
	SourceFormat string                 `json:"source_format"`
	TargetFormat string                 `json:"target_format"`
	CreatedAt    time.Time              `json:"created_at"`
}

// UploadFile handles file upload and creates a conversion job
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Get user from authenticated context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// Extract user ID from user object
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Message: "Invalid user data",
		})
		return
	}

	userID := userModel.ID

	// Parse form data
	var uploadReq UploadRequest
	if err := c.ShouldBind(&uploadReq); err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Invalid form data: " + err.Error(),
		})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "No file provided",
		})
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: fmt.Sprintf("File too large. Maximum size: %dMB", MaxFileSize/1024/1024),
		})
		return
	}

	// Validate file type
	originalFilename := header.Filename
	sourceFormat := strings.ToLower(filepath.Ext(originalFilename))
	if sourceFormat == "" {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Cannot determine file type from filename",
		})
		return
	}

	// Remove leading dot from extension
	sourceFormat = strings.TrimPrefix(sourceFormat, ".")

	// Validate target format
	targetFormat := uploadReq.TargetFormat
	if targetFormat == "" {
		targetFormat = sourceFormat // Default to same format
	}

	// Validate file category
	category, err := storage.ValidateFileType(originalFilename)
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Save file locally
	inputPath, err := h.localStorage.SaveFile(file, originalFilename, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Failed to save file: " + err.Error(),
		})
		return
	}

	// Generate job ID
	jobID := uuid.New().String()

	// Create job using job service (which adds to Redis queue)
	job := models.Job{
		JobID:            jobID,
		UserID:           userID,
		OriginalFilename: originalFilename,
		FileSize:         header.Size,
		SourceFormat:     sourceFormat,
		TargetFormat:     targetFormat,
		Status:           string(models.StatusPending),
		InputPath:        inputPath,
		OutputPath:       "", // Will be set when processing completes
		Error:            "",
	}

	// Use job service to create job (adds to database AND Redis queue)
	if h.jobService != nil {
		ctx := context.Background()
		settings := map[string]interface{}{
			"quality_preset": uploadReq.QualityPreset,
		}
		if err := h.jobService.CreateJob(ctx, &job, settings); err != nil {
			// Clean up uploaded file on error
			h.localStorage.DeleteFile(inputPath)
			c.JSON(http.StatusInternalServerError, UploadResponse{
				Success: false,
				Message: "Failed to create job: " + err.Error(),
			})
			return
		}
	} else {
		// Fallback: create job directly in database (no Redis queue)
		if err := h.db.Create(&job).Error; err != nil {
			// Clean up uploaded file on error
			h.localStorage.DeleteFile(inputPath)
			c.JSON(http.StatusInternalServerError, UploadResponse{
				Success: false,
				Message: "Failed to create job: " + err.Error(),
			})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusCreated, UploadResponse{
		Success:      true,
		Message:      fmt.Sprintf("%s conversion job created successfully", category),
		JobID:        jobID,
		Status:       string(models.StatusPending),
		OriginalName: originalFilename,
		FileSize:     header.Size,
		SourceFormat: sourceFormat,
		TargetFormat: targetFormat,
		CreatedAt:    job.CreatedAt,
		FileInfo: map[string]interface{}{
			"category":       category,
			"quality_preset": uploadReq.QualityPreset,
		},
	})
}

// GetUserJobs returns all jobs for the authenticated user
func (h *UploadHandler) GetUserJobs(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Extract user ID from user object
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user data",
		})
		return
	}

	userID := userModel.ID

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var jobs []models.Job
	var total int64

	// Get total count
	h.db.Model(&models.Job{}).Where("user_id = ?", userID).Count(&total)

	// Get paginated results
	result := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch jobs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetJobStatus returns the status of a specific job
func (h *UploadHandler) GetJobStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Extract user ID from user object
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user data",
		})
		return
	}

	userID := userModel.ID

	jobID := c.Param("id")

	var job models.Job
	result := h.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&job)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch job",
		})
		return
	}

	// Get download URL if job is completed
	var downloadURL string
	if job.Status == string(models.StatusCompleted) && job.OutputPath != "" {
		downloadURL = fmt.Sprintf("/api/v1/download/%s", jobID)
	}

	c.JSON(http.StatusOK, gin.H{
		"job_id":            job.JobID,
		"status":            job.Status,
		"original_filename": job.OriginalFilename,
		"file_size":         job.FileSize,
		"source_format":     job.SourceFormat,
		"target_format":     job.TargetFormat,
		"input_path":        job.InputPath,
		"output_path":       job.OutputPath,
		"error":             job.Error,
		"created_at":        job.CreatedAt,
		"updated_at":        job.UpdatedAt,
		"download_url":      downloadURL,
	})
}

// DownloadFile serves the processed file for download
func (h *UploadHandler) DownloadFile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Extract user ID from user object
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user data",
		})
		return
	}

	userID := userModel.ID

	jobID := c.Param("id")

	var job models.Job
	result := h.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&job)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch job",
		})
		return
	}

	// Check if job is completed
	if job.Status != string(models.StatusCompleted) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job not completed",
		})
		return
	}

	// Check if output file exists
	if job.OutputPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Output file path not set",
		})
		return
	}

	// Serve the file
	c.FileAttachment(job.OutputPath, job.OriginalFilename)
}
