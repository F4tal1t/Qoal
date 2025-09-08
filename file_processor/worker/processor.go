package worker

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"

	"github.com/go-redis/redis/v8"
)

type Processor struct {
	redisClient      *redis.Client
	jobService       *services.JobService
	imageProcessor   *services.EnhancedImageProcessor
	audioProcessor   *services.EnhancedAudioProcessor
	archiveProcessor *services.ArchiveProcessor
}

func NewProcessor(
	redisClient *redis.Client,
	s3Client interface{},
	jobService *services.JobService,
	imageProcessor *services.EnhancedImageProcessor,
	audioProcessor *services.EnhancedAudioProcessor,
	archiveProcessor *services.ArchiveProcessor,
) *Processor {
	return &Processor{
		redisClient:      redisClient,
		jobService:       jobService,
		imageProcessor:   imageProcessor,
		audioProcessor:   audioProcessor,
		archiveProcessor: archiveProcessor,
	}
}

func (p *Processor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
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

	// Create processing job model
	processingJob := &models.ProcessingJob{
		JobID:        job.ID,
		InputPath:    job.InputFile,
		SourceFormat: strings.TrimPrefix(filepath.Ext(job.InputFile), "."),
		Status:       models.StatusPending,
	}

	// Process based on file type
	switch {
	case strings.HasSuffix(job.InputFile, ".jpg") || strings.HasSuffix(job.InputFile, ".png"):
		err = p.imageProcessor.ProcessImage(processingJob)
	case strings.HasSuffix(job.InputFile, ".mp3") || strings.HasSuffix(job.InputFile, ".wav"):
		err = p.audioProcessor.ProcessAudio(processingJob)
	case strings.HasSuffix(job.InputFile, ".zip") || strings.HasSuffix(job.InputFile, ".tar"):
		err = p.archiveProcessor.ProcessArchive(processingJob)
	default:
		return fmt.Errorf("unsupported file type")
	}

	if err != nil {
		return err
	}

	// Update job status
	_, err = p.jobService.UpdateJobStatus(ctx, job.ID, processingJob.Status, processingJob.OutputPath)
	return err
}
