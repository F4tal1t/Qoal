package handlers

import (
	"context"
	"net/http"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) CreateJobHandler(c *gin.Context) {
	// Get user ID from authenticated context
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get user model and extract ID
	userModel, ok := userID.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user data",
		})
		return
	}
	userIDStr := userModel.ID

	// Parse request body
	var req struct {
		InputPath    string                 `json:"input_path" binding:"required"`
		OutputPath   string                 `json:"output_path"`
		SourceFormat string                 `json:"source_format" binding:"required"`
		TargetFormat string                 `json:"target_format" binding:"required"`
		Settings     map[string]interface{} `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Generate job ID
	jobID := uuid.New().String()

	// Create job
	job := models.Job{
		JobID:        jobID,
		UserID:       userIDStr,
		InputPath:    req.InputPath,
		OutputPath:   req.OutputPath,
		SourceFormat: req.SourceFormat,
		TargetFormat: req.TargetFormat,
		Status:       string(models.StatusPending),
	}

	// Add job to processing queue
	ctx := context.Background()
	if err := h.jobService.CreateJob(ctx, &job, req.Settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":     jobID,
		"status":     "job_created",
		"message":    "Job queued for processing",
		"created_at": job.CreatedAt,
	})
}

func (h *JobHandler) GetJobStatusHandler(c *gin.Context) {
	// Get user ID from authenticated context
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get user model and extract ID
	userModel, ok := userID.(*models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user data",
		})
		return
	}
	userIDStr := userModel.ID

	jobID := c.Param("id")

	ctx := context.Background()
	job, err := h.jobService.GetJob(ctx, jobID, userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
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
	})
}
