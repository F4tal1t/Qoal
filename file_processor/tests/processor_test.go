package tests

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/qoal/file-processor/models"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/worker"
)

type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) UpdateJob(job *models.ProcessingJob) error {
	args := m.Called(job)
	return args.Error(0)
}

type MockS3Service struct {
	mock.Mock
}

func (m *MockS3Service) DownloadFile(sourcePath, destinationPath string) error {
	args := m.Called(sourcePath, destinationPath)
	return args.Error(0)
}

func (m *MockS3Service) UploadFile(sourcePath, destinationPath string) error {
	args := m.Called(sourcePath, destinationPath)
	return args.Error(0)
}

func TestProcessor_ProcessJob(t *testing.T) {
	// Setup
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	s3Client := &s3.S3{}
	mockJobService := new(MockJobService)
	mockS3Service := new(MockS3Service)

	config := &config.Config{
		ImageMagickPath: "/usr/bin/convert",
		FFmpegPath:      "/usr/bin/ffmpeg",
		TempDir:         os.TempDir(),
		OutputDir:       "./test_outputs",
	}

	imageProcessor := services.NewEnhancedImageProcessor(config, mockS3Service)
	audioProcessor := services.NewEnhancedAudioProcessor(config, mockS3Service)
	archiveProcessor := services.NewArchiveProcessor(config, mockS3Service)

	processor := worker.NewProcessor(
		redisClient,
		s3Client,
		mockJobService,
		imageProcessor,
		audioProcessor,
		archiveProcessor,
	)

	t.Run("Image Conversion", func(t *testing.T) {
		job := &models.ProcessingJob{
			Type:         "image_conversion",
			SourceFormat: "jpg",
			TargetFormat: "png",
		}

		mockS3Service.On("DownloadFile", mock.Anything, mock.Anything).Return(nil)
		mockS3Service.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
		mockJobService.On("UpdateJob", mock.Anything).Return(nil)

		err := processor.processJob(job)
		assert.NoError(t, err)
		mockS3Service.AssertExpectations(t)
		mockJobService.AssertExpectations(t)
	})

	t.Run("Audio Conversion", func(t *testing.T) {
		job := &models.ProcessingJob{
			Type:         "audio_conversion",
			SourceFormat: "mp3",
			TargetFormat: "wav",
		}

		mockS3Service.On("DownloadFile", mock.Anything, mock.Anything).Return(nil)
		mockS3Service.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
		mockJobService.On("UpdateJob", mock.Anything).Return(nil)

		err := processor.processJob(job)
		assert.NoError(t, err)
		mockS3Service.AssertExpectations(t)
		mockJobService.AssertExpectations(t)
	})

	t.Run("Archive Conversion", func(t *testing.T) {
		job := &models.ProcessingJob{
			Type:         "archive_conversion",
			SourceFormat: "zip",
			TargetFormat: "7z",
		}

		mockS3Service.On("DownloadFile", mock.Anything, mock.Anything).Return(nil)
		mockS3Service.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
		mockJobService.On("UpdateJob", mock.Anything).Return(nil)

		err := processor.processJob(job)
		assert.NoError(t, err)
		mockS3Service.AssertExpectations(t)
		mockJobService.AssertExpectations(t)
	})
}
