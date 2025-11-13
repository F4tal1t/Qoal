package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/storage"
)

type ProcessorS3 struct {
	redisClient       *redis.Client
	jobService        *services.JobService
	config            *config.Config
	s3Storage         *storage.S3Storage
	documentProcessor *services.EnhancedDocumentProcessor
	imageProcessor    *services.EnhancedImageProcessor
	videoProcessor    *services.EnhancedVideoProcessor
	audioProcessor    *services.EnhancedAudioProcessor
	archiveProcessor  *services.ArchiveProcessor
}

func NewProcessorS3(jobService *services.JobService, cfg *config.Config, redisClient *redis.Client, s3Storage *storage.S3Storage) *ProcessorS3 {
	return &ProcessorS3{
		redisClient:       redisClient,
		jobService:        jobService,
		config:            cfg,
		s3Storage:         s3Storage,
		documentProcessor: services.NewEnhancedDocumentProcessor(cfg),
		imageProcessor:    services.NewEnhancedImageProcessor(cfg),
		videoProcessor:    services.NewEnhancedVideoProcessor(cfg),
		audioProcessor:    services.NewEnhancedAudioProcessor(cfg),
		archiveProcessor:  services.NewArchiveProcessor(cfg),
	}
}

func (p *ProcessorS3) Start(ctx context.Context) {
	log.Println("Starting S3 job processor worker...")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("S3 Worker shutting down...")
				return
			default:
				task, err := p.jobService.GetNextJobFromQueue(ctx)
				if err != nil {
					if err != redis.Nil {
						log.Printf("Error getting job from queue: %v", err)
					}
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

func (p *ProcessorS3) ProcessJob(ctx context.Context, task *services.JobTask) error {
	if err := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusProcessing, "", ""); err != nil {
		return fmt.Errorf("failed to update job status to processing: %w", err)
	}

	// Download from S3 to temp
	tempInput := filepath.Join(p.config.TempDir, task.JobID+"_input"+filepath.Ext(task.InputPath))
	inputFile, err := os.Create(tempInput)
	if err != nil {
		return fmt.Errorf("failed to create temp input file: %w", err)
	}

	s3File, err := p.s3Storage.GetFile(task.InputPath)
	if err != nil {
		inputFile.Close()
		os.Remove(tempInput)
		return fmt.Errorf("failed to download from S3: %w", err)
	}

	_, err = io.Copy(inputFile, s3File)
	s3File.Close()
	inputFile.Close()
	if err != nil {
		os.Remove(tempInput)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	processingJob := &models.ProcessingJob{
		JobID:        task.JobID,
		UserID:       task.UserID,
		InputPath:    tempInput,
		OutputPath:   "",
		SourceFormat: task.SourceFormat,
		TargetFormat: task.TargetFormat,
		Status:       models.StatusProcessing,
		Settings:     task.Settings,
	}

	log.Printf("Processing file: %s (format: %s -> %s)", task.InputPath, task.SourceFormat, task.TargetFormat)

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

	os.Remove(tempInput)

	if err != nil {
		if updateErr := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusFailed, "", err.Error()); updateErr != nil {
			log.Printf("Failed to update job status to failed: %v", updateErr)
		}
		return fmt.Errorf("job processing failed: %w", err)
	}

	// Upload result to S3
	outputFile, err := os.Open(processingJob.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer outputFile.Close()
	defer os.Remove(processingJob.OutputPath)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, outputFile)
	if err != nil {
		return fmt.Errorf("failed to read output file: %w", err)
	}

	s3OutputPath, err := p.s3Storage.SaveProcessedFile(bytes.NewReader(buf.Bytes()), task.JobID, task.TargetFormat)
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	if err := p.jobService.UpdateJobStatus(ctx, task.JobID, models.StatusCompleted, s3OutputPath, ""); err != nil {
		return fmt.Errorf("failed to update job status to completed: %w", err)
	}

	log.Printf("Job %s completed successfully", task.JobID)
	return nil
}

func (p *ProcessorS3) getFileCategory(format string) string {
	format = strings.ToLower(format)

	documentFormats := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf", "csv"}
	for _, f := range documentFormats {
		if format == f {
			return "document"
		}
	}

	imageFormats := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "tiff", "svg"}
	for _, f := range imageFormats {
		if format == f {
			return "image"
		}
	}

	videoFormats := []string{"mp4", "avi", "mov", "wmv", "flv", "mkv", "webm", "m4v"}
	for _, f := range videoFormats {
		if format == f {
			return "video"
		}
	}

	audioFormats := []string{"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma"}
	for _, f := range audioFormats {
		if format == f {
			return "audio"
		}
	}

	archiveFormats := []string{"zip", "rar", "7z", "tar", "gz", "bz2", "xz"}
	for _, f := range archiveFormats {
		if format == f {
			return "archive"
		}
	}

	return "unknown"
}
