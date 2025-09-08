package worker

import (
	"context"
	"log"
	"time"

	"github.com/qoal/file-processor/services"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis/v8"
)

type Processor struct {
	redisClient *redis.Client
	s3Client    *s3.S3
	jobService  *services.JobService

	// Processors
	imageProcessor   *services.EnhancedImageProcessor
	audioProcessor   *services.EnhancedAudioProcessor
	archiveProcessor *services.ArchiveProcessor
}

func NewProcessor(
	redisClient *redis.Client,
	s3Client *s3.S3,
	jobService *services.JobService,
	imageProcessor *services.EnhancedImageProcessor,
	audioProcessor *services.EnhancedAudioProcessor,
	archiveProcessor *services.ArchiveProcessor,
) *Processor {
	return &Processor{
		redisClient: redisClient,
		s3Client:    s3Client,
		jobService:  jobService,

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
				// TODO: Implement job processing logic
				log.Println("Checking for jobs...")
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
