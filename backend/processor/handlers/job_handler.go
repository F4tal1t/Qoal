package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jobService *services.JobService

func SetJobService(js *services.JobService) {
	jobService = js
}

func CreateJobHandler(c *gin.Context) {
	jobID := uuid.New().String()

	job := models.Job{
		ID:        jobID,
		UserID:    c.PostForm("user_id"),
		InputFile: c.PostForm("input_file"),
		Status:    models.StatusPending,
		CreatedAt: time.Now().Unix(),
	}

	// Add job to processing queue
	ctx := context.Background()
	if err := jobService.CreateJob(ctx, &job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "job_created",
	})
}

func GetJobStatusHandler(c *gin.Context) {
	jobID := c.Param("id")

	ctx := context.Background()
	job, err := jobService.GetJobStatus(ctx, jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_id":      job.ID,
		"status":      job.Status,
		"input_file":  job.InputFile,
		"output_file": job.OutputFile,
		"created_at":  job.CreatedAt,
		"completed_at": job.CompletedAt,
	})
}
