package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/storage"

	"github.com/go-redis/redis/v8"
)

type Processor struct {
	redisClient       *redis.Client
	jobService        *services.JobService
	config            *config.Config
	localStorage      *storage.LocalStorage
	documentProcessor *services.EnhancedDocumentProcessor
	imageProcessor    *services.EnhancedImageProcessor
	videoProcessor    *services.EnhancedVideoProcessor
	audioProcessor    *services.EnhancedAudioProcessor
	archiveProcessor  *services.ArchiveProcessor
}

func NewProcessor(jobService *services.JobService, cfg *config.Config, redisClient *redis.Client, localStorage *storage.LocalStorage) *Processor {
	return &Processor{
		redisClient:       redisClient,
		jobService:        jobService,
		config:            cfg,
		localStorage:      localStorage,
		documentProcessor: services.NewEnhancedDocumentProcessor(cfg, localStorage),
		imageProcessor:    services.NewEnhancedImageProcessor(cfg, localStorage),
		videoProcessor:    services.NewEnhancedVideoProcessor(cfg, localStorage),
		audioProcessor:    services.NewEnhancedAudioProcessor(cfg, localStorage),
		archiveProcessor:  services.NewArchiveProcessor(cfg, localStorage),
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
				// Get next job from queue
				task, err := p.jobService.GetNextJobFromQueue(ctx)
				if err != nil {
					if err != redis.Nil {
						log.Printf("Error getting job from queue: %v", err)
					}
					// Wait a bit before trying again
					time.Sleep(1 * time.Second)
					continue
				}

				log.Printf("Processing job: %s", task.JobID)
				if err := p.ProcessJob(ctx, task); err != nil {
					log.Printf("Job processing failed: %v", err)
				}
			}
		}
	}()
}

func (p *Processor) ProcessJob(ctx context.Context, task *services.JobTask) error {
	// Update job status to processing
	if err := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusProcessing, "", ""); err != nil {
		return fmt.Errorf("failed to update job status to processing: %w", err)
	}

	// Create processing job model
	processingJob := &models.ProcessingJob{
		JobID:        task.JobID,
		UserID:       task.UserID,
		InputPath:    task.InputPath,
		OutputPath:   task.OutputPath,
		SourceFormat: task.SourceFormat,
		TargetFormat: task.TargetFormat,
		Status:       models.StatusProcessing,
		Settings:     task.Settings,
	}

	log.Printf("Processing file: %s (format: %s -> %s)", task.InputPath, task.SourceFormat, task.TargetFormat)

	// Determine file category and process accordingly
	var err error
	fileCategory := p.getFileCategory(task.SourceFormat)

	switch fileCategory {
	case "document":
		err = p.documentProcessor.ProcessDocument(processingJob)
	case "image":
		err = p.imageProcessor.ProcessImage(processingJob)
	case "video":
		err = p.videoProcessor.ProcessVideo(processingJob)
	case "audio":
		err = p.audioProcessor.ProcessAudio(processingJob)
	case "archive":
		err = p.archiveProcessor.ProcessArchive(processingJob)
	default:
		err = fmt.Errorf("unsupported file category: %s", fileCategory)
	}

	if err != nil {
		// Update job status to failed
		if updateErr := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusFailed, "", err.Error()); updateErr != nil {
			log.Printf("Failed to update job status to failed: %v", updateErr)
		}
		return fmt.Errorf("job processing failed: %w", err)
	}

	// Update job status to completed
	if err := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusCompleted, processingJob.OutputPath, ""); err != nil {
		return fmt.Errorf("failed to update job status to completed: %w", err)
	}

	log.Printf("Job %s completed successfully", task.JobID)
	return nil
}

// getFileCategory determines the file category based on format
func (p *Processor) getFileCategory(format string) string {
	format = strings.ToLower(format)

	// Document formats
	documentFormats := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf", "csv"}
	for _, f := range documentFormats {
		if format == f {
			return "document"
		}
	}

	// Image formats
	imageFormats := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "tiff", "svg"}
	for _, f := range imageFormats {
		if format == f {
			return "image"
		}
	}

	// Video formats
	videoFormats := []string{"mp4", "avi", "mov", "wmv", "flv", "mkv", "webm", "m4v"}
	for _, f := range videoFormats {
		if format == f {
			return "video"
		}
	}

	// Audio formats
	audioFormats := []string{"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma"}
	for _, f := range audioFormats {
		if format == f {
			return "audio"
		}
	}

	// Archive formats
	archiveFormats := []string{"zip", "rar", "7z", "tar", "gz", "bz2", "xz"}
	for _, f := range archiveFormats {
		if format == f {
			return "archive"
		}
	}

	return "unknown"
}
