package services

import (
	"context"
	"fmt"
	"time"

	"github.com/qoal/file-processor/models"

	"github.com/go-redis/redis/v8"
)

type JobService struct {
	redisClient *redis.Client
}

func NewJobService(redisClient *redis.Client) *JobService {
	return &JobService{
		redisClient: redisClient,
	}
}

func (js *JobService) CreateJob(ctx context.Context, job *models.Job) error {
	// Store job in Redis with expiration
	jobKey := "job:" + job.ID
	jobData := map[string]interface{}{
		"id":         job.ID,
		"user_id":    job.UserID,
		"input_file": job.InputFile,
		"status":     string(job.Status),
		"created_at": job.CreatedAt,
	}

	// Store job data
	if err := js.redisClient.HMSet(ctx, jobKey, jobData).Err(); err != nil {
		return err
	}

	// Set expiration (24 hours)
	js.redisClient.Expire(ctx, jobKey, 24*time.Hour)

	// Add to processing queue
	return js.redisClient.LPush(ctx, "processing_queue", job.ID).Err()
}

func (js *JobService) GetJobStatus(ctx context.Context, jobID string) (*models.Job, error) {
	jobKey := "job:" + jobID
	jobData, err := js.redisClient.HGetAll(ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	if len(jobData) == 0 {
		return nil, fmt.Errorf("job not found")
	}

	job := &models.Job{
		ID:        jobData["id"],
		UserID:    jobData["user_id"],
		InputFile: jobData["input_file"],
		Status:    models.JobStatus(jobData["status"]),
	}

	if outputFile, exists := jobData["output_file"]; exists {
		job.OutputFile = outputFile
	}

	return job, nil
}

func (js *JobService) UpdateJobStatus(ctx context.Context, jobID string, status models.JobStatus, outputPath string) (*models.Job, error) {
	jobKey := "job:" + jobID
	updateData := map[string]interface{}{
		"status": string(status),
	}

	if outputPath != "" {
		updateData["output_file"] = outputPath
	}

	if status == models.StatusCompleted || status == models.StatusFailed {
		updateData["completed_at"] = time.Now().Unix()
	}

	if err := js.redisClient.HMSet(ctx, jobKey, updateData).Err(); err != nil {
		return nil, err
	}

	return js.GetJobStatus(ctx, jobID)
}
