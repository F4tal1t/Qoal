package handlers

import (
	"net/http"
	"time"

	"github.com/qoal/file-processor/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateJobHandler(c *gin.Context) {
	jobID := uuid.New().String()

	_ = models.Job{
		ID:        jobID,
		InputFile: c.PostForm("file"),
		Status:    models.StatusPending,
		CreatedAt: time.Now().Unix(),
	}

	// TODO: Add job to processing queue

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "job_created",
	})
}

func GetJobStatusHandler(c *gin.Context) {
	jobID := c.Param("id")
	// TODO: Get job status from database/queue

	c.JSON(http.StatusOK, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}
