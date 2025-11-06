package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/qoal/file-processor/models"
)

type JobService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewJobService(db *gorm.DB, redisClient *redis.Client) *JobService {
	return &JobService{
		db:          db,
		redisClient: redisClient,
	}
}

type JobTask struct {
	JobID        string                 `json:"job_id"`
	UserID       string                 `json:"user_id"`
	InputPath    string                 `json:"input_path"`
	OutputPath   string                 `json:"output_path"`
	SourceFormat string                 `json:"source_format"`
	TargetFormat string                 `json:"target_format"`
	Settings     map[string]interface{} `json:"settings"`
	CreatedAt    time.Time              `json:"created_at"`
}

// CreateJob creates a new processing job
func (s *JobService) CreateJob(ctx context.Context, job *models.Job, settings map[string]interface{}) error {
	// Create job in database
	if err := s.db.Create(job).Error; err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	// Create job task for queue
	task := JobTask{
		JobID:        job.JobID,
		UserID:       job.UserID,
		InputPath:    job.InputPath,
		OutputPath:   job.OutputPath,
		SourceFormat: job.SourceFormat,
		TargetFormat: job.TargetFormat,
		Settings:     settings,
		CreatedAt:    job.CreatedAt,
	}

	// Add job to Redis queue
	taskData, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal job task: %w", err)
	}

	if err := s.redisClient.LPush(ctx, "conversion_queue", taskData).Err(); err != nil {
		return fmt.Errorf("failed to add job to queue: %w", err)
	}

	return nil
}

// GetJob retrieves a job by ID
func (s *JobService) GetJob(ctx context.Context, jobID string, userID string) (*models.Job, error) {
	var job models.Job
	if err := s.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return &job, nil
}

// UpdateJobStatus updates the status of a job
func (s *JobService) UpdateJobStatus(ctx context.Context, jobID string, status models.JobStatus, outputPath string, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     string(status),
		"updated_at": time.Now(),
	}

	if outputPath != "" {
		updates["output_path"] = outputPath
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	if status == models.StatusCompleted {
		updates["completed_at"] = time.Now()
	}

	if err := s.db.Model(&models.Job{}).Where("job_id = ?", jobID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// GetUserJobs retrieves all jobs for a user with pagination
func (s *JobService) GetUserJobs(ctx context.Context, userID string, page, limit int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	// Get total count
	if err := s.db.Model(&models.Job{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get jobs: %w", err)
	}

	return jobs, total, nil
}

// GetNextJobFromQueue retrieves the next job from the Redis queue
func (s *JobService) GetNextJobFromQueue(ctx context.Context) (*JobTask, error) {
	result, err := s.redisClient.BRPop(ctx, 0, "conversion_queue").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get job from queue: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result format")
	}

	var task JobTask
	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job task: %w", err)
	}

	return &task, nil
}

// DeleteJob deletes a job and its associated files
func (s *JobService) DeleteJob(ctx context.Context, jobID string, userID string) error {
	var job models.Job
	if err := s.db.Where("job_id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("job not found")
		}
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Delete from database
	if err := s.db.Delete(&job).Error; err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	return nil
}
