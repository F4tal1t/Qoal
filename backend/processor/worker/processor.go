package worker

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"

	"github.com/go-redis/redis/v8"
)

type Processor struct {
	redisClient *redis.Client
	jobService  *services.JobService
	config      *config.Config
}

func NewProcessor(jobService *services.JobService, cfg *config.Config, redisClient *redis.Client) *Processor {
	return &Processor{
		redisClient: redisClient,
		jobService:  jobService,
		config:      cfg,
	}
}

func (p *Processor) Start(ctx context.Context) {
	log.Println("Starting job processor worker...")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Worker shutting down...")
				return
			default:
				// Check for jobs in Redis queue
				jobID, err := p.redisClient.BRPop(ctx, 5*time.Second, "processing_queue").Result()
				if err != nil {
					if err != redis.Nil {
						log.Printf("Error checking queue: %v", err)
					}
					continue
				}

				if len(jobID) > 1 {
					log.Printf("Processing job: %s", jobID[1])
					if err := p.ProcessJob(ctx, jobID[1]); err != nil {
						log.Printf("Job processing failed: %v", err)
					}
				}
			}
		}
	}()
}

func (p *Processor) ProcessJob(ctx context.Context, jobID string) error {
	job, err := p.jobService.GetJobStatus(ctx, jobID)
	if err != nil {
		return err
	}

	// Update job status to processing
	_, err = p.jobService.UpdateJobStatus(ctx, job.ID, models.StatusProcessing, "")
	if err != nil {
		return err
	}

	// Create processing job model
	processingJob := &models.ProcessingJob{
		JobID:        job.ID,
		InputPath:    job.InputFile,
		SourceFormat: strings.TrimPrefix(filepath.Ext(job.InputFile), "."),
		Status:       models.StatusProcessing,
	}

	// For now, just simulate processing and mark as completed
	// TODO: Implement actual file processing based on file type
	log.Printf("Simulating processing for file: %s", job.InputFile)

	// Simulate processing time
	time.Sleep(2 * time.Second)

	// Mark as completed
	processingJob.Status = models.StatusCompleted
	processingJob.OutputPath = job.InputFile + ".processed" // Placeholder

	// Update job status
	_, err = p.jobService.UpdateJobStatus(ctx, job.ID, processingJob.Status, processingJob.OutputPath)
	if err != nil {
		return err
	}

	log.Printf("Job %s completed successfully", jobID)
	return nil
}