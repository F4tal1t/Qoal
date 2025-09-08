package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qoal/file-processor/handlers"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/worker"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Initialize services
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))

	s3Client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	// Create Gin router
	r := gin.Default()

	// Initialize services
	config := &config.Config{
		ImageMagickPath: os.Getenv("IMAGEMAGICK_PATH"),
		FFmpegPath:      os.Getenv("FFMPEG_PATH"),
		TempDir:         os.TempDir(),
		OutputDir:       "./outputs",
	}

	s3Service := services.NewS3Service(downloader, uploader, os.Getenv("AWS_S3_BUCKET"))
	jobService := services.NewJobService(rdb)

	// Initialize processors
	imageProcessor := services.NewEnhancedImageProcessor(config, s3Service)
	audioProcessor := services.NewEnhancedAudioProcessor(config, s3Service)
	archiveProcessor := services.NewArchiveProcessor(config, s3Service)

	// Initialize worker with all processors
	worker := worker.NewProcessor(rdb, s3Client, jobService, imageProcessor, audioProcessor, archiveProcessor)
	worker.Start(context.Background())

	// Setup routes
	r.POST("/process", handlers.CreateJobHandler)
	r.GET("/status/:id", handlers.GetJobStatusHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
