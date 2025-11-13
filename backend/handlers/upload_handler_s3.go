package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/storage"
)

type UploadHandlerS3 struct {
	db         *gorm.DB
	s3Storage  *storage.S3Storage
	jobService *services.JobService
}

func NewUploadHandlerS3(db *gorm.DB, s3Storage *storage.S3Storage, jobService *services.JobService) *UploadHandler {
	return &UploadHandler{
		db:           db,
		localStorage: nil,
		jobService:   jobService,
	}
}

func (h *UploadHandler) UploadFileS3(c *gin.Context, s3Storage *storage.S3Storage) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, UploadResponse{Success: false, Message: "User not authenticated"})
		return
	}

	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, UploadResponse{Success: false, Message: "Invalid user data"})
		return
	}

	var uploadReq UploadRequest
	if err := c.ShouldBind(&uploadReq); err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{Success: false, Message: "Invalid form data: " + err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{Success: false, Message: "No file provided"})
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, UploadResponse{Success: false, Message: fmt.Sprintf("File too large. Maximum size: %dMB", MaxFileSize/1024/1024)})
		return
	}

	originalFilename := header.Filename
	sourceFormat := strings.ToLower(filepath.Ext(originalFilename))
	if sourceFormat == "" {
		c.JSON(http.StatusBadRequest, UploadResponse{Success: false, Message: "Cannot determine file type from filename"})
		return
	}

	sourceFormat = strings.TrimPrefix(sourceFormat, ".")
	targetFormat := uploadReq.TargetFormat
	if targetFormat == "" {
		targetFormat = sourceFormat
	}

	category, err := storage.ValidateFileType(originalFilename)
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadResponse{Success: false, Message: err.Error()})
		return
	}

	inputPath, err := s3Storage.SaveFile(file, originalFilename, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadResponse{Success: false, Message: "Failed to save file: " + err.Error()})
		return
	}

	jobID := uuid.New().String()
	job := models.Job{
		JobID:            jobID,
		UserID:           userModel.ID,
		OriginalFilename: originalFilename,
		FileSize:         header.Size,
		SourceFormat:     sourceFormat,
		TargetFormat:     targetFormat,
		Status:           string(models.StatusPending),
		InputPath:        inputPath,
		OutputPath:       "",
		Error:            "",
	}

	if h.jobService != nil {
		ctx := context.Background()
		settings := map[string]interface{}{"quality_preset": uploadReq.QualityPreset}
		if err := h.jobService.CreateJob(ctx, &job, settings); err != nil {
			s3Storage.DeleteFile(inputPath)
			c.JSON(http.StatusInternalServerError, UploadResponse{Success: false, Message: "Failed to create job: " + err.Error()})
			return
		}
	} else {
		if err := h.db.Create(&job).Error; err != nil {
			s3Storage.DeleteFile(inputPath)
			c.JSON(http.StatusInternalServerError, UploadResponse{Success: false, Message: "Failed to create job: " + err.Error()})
			return
		}
	}

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
		FileInfo:     map[string]interface{}{"category": category, "quality_preset": uploadReq.QualityPreset},
	})
}

func (h *UploadHandler) DownloadFileS3(c *gin.Context, s3Storage *storage.S3Storage) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
		return
	}

	jobID := c.Param("id")
	var job models.Job
	result := h.db.Where("job_id = ? AND user_id = ?", jobID, userModel.ID).First(&job)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}

	if job.Status != string(models.StatusCompleted) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not completed"})
		return
	}

	if job.OutputPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Output file path not set"})
		return
	}

	presignedURL, err := s3Storage.GetPresignedURL(job.OutputPath, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": presignedURL,
		"filename":     job.OriginalFilename,
		"job_id":       job.JobID,
	})
}
