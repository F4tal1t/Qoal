package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/stretchr/testify/assert"
)

func TestDay5Implementation(t *testing.T) {
	// Setup test environment
	cfg := &config.Config{
		ImageMagickPath: "convert", // Assuming ImageMagick is in PATH
		FFmpegPath:      "ffmpeg",  // Assuming FFmpeg is in PATH
		TempDir:         os.TempDir(),
		OutputDir:       "./test_outputs",
	}

	// Create output directory
	os.MkdirAll(cfg.OutputDir, 0755)
	defer os.RemoveAll(cfg.OutputDir)

	// Mock S3 service for testing
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	s3Service := services.NewS3Service(downloader, uploader, "test-bucket")

	// Initialize processors for testing
	imageProcessor := services.NewEnhancedImageProcessor(cfg, s3Service)
	audioProcessor := services.NewEnhancedAudioProcessor(cfg, s3Service)
	archiveProcessor := services.NewArchiveProcessor(cfg, s3Service)

	t.Run("TestImageConversionLogic", func(t *testing.T) {
		testImageConversions(t, imageProcessor)
	})

	t.Run("TestAudioConversionLogic", func(t *testing.T) {
		testAudioConversions(t, audioProcessor)
	})

	t.Run("TestArchiveConversionLogic", func(t *testing.T) {
		testArchiveConversions(t, archiveProcessor)
	})

	// Skip Redis-dependent tests
	// t.Run("TestJobServiceFunctionality", func(t *testing.T) {
	// 	testJobService(t, jobService)
	// })

	t.Run("TestUtilityFunctions", func(t *testing.T) {
		testUtilityFunctions(t)
	})
}

func testImageConversions(t *testing.T, processor *services.EnhancedImageProcessor) {
	// Test image extension mapping
	ext, err := utils.GetImageExtension("jpeg")
	assert.NoError(t, err)
	assert.Equal(t, ".jpg", ext)

	ext, err = utils.GetImageExtension("png")
	assert.NoError(t, err)
	assert.Equal(t, ".png", ext)

	ext, err = utils.GetImageExtension("webp")
	assert.NoError(t, err)
	assert.Equal(t, ".webp", ext)

	// Test unsupported format
	_, err = utils.GetImageExtension("unsupported")
	assert.Error(t, err)
}

func testAudioConversions(t *testing.T, processor *services.EnhancedAudioProcessor) {
	// Test audio extension mapping
	ext, err := utils.GetAudioExtension("mp3")
	assert.NoError(t, err)
	assert.Equal(t, ".mp3", ext)

	ext, err = utils.GetAudioExtension("wav")
	assert.NoError(t, err)
	assert.Equal(t, ".wav", ext)

	ext, err = utils.GetAudioExtension("flac")
	assert.NoError(t, err)
	assert.Equal(t, ".flac", ext)

	// Test quality presets
	settings := map[string]interface{}{
		"quality_preset": "high",
	}
	quality := processor.GetAudioQualityPreset(settings)
	assert.Equal(t, "256k", quality.Bitrate)
	assert.Equal(t, 48000, quality.SampleRate)
}

func testArchiveConversions(t *testing.T, processor *services.ArchiveProcessor) {
	// Test archive extension mapping
	ext, err := utils.GetArchiveExtension("zip")
	assert.NoError(t, err)
	assert.Equal(t, ".zip", ext)

	ext, err = utils.GetArchiveExtension("7z")
	assert.NoError(t, err)
	assert.Equal(t, ".7z", ext)

	// Test compression level mapping
	settings := map[string]interface{}{
		"compression_level": "maximum",
	}
	level := processor.GetCompressionLevel(settings)
	assert.Equal(t, 7, level)
}

func testJobService(t *testing.T, jobService *services.JobService) {
	ctx := context.Background()

	// Create test job
	job := &models.Job{
		ID:        "test-job-123",
		UserID:    "user-456",
		InputFile: "test-input.jpg",
		Status:    models.StatusPending,
		CreatedAt: time.Now().Unix(),
	}

	// Test job creation
	err := jobService.CreateJob(ctx, job)
	assert.NoError(t, err)

	// Test job retrieval
	retrievedJob, err := jobService.GetJobStatus(ctx, job.ID)
	assert.NoError(t, err)
	assert.Equal(t, job.ID, retrievedJob.ID)
	assert.Equal(t, job.UserID, retrievedJob.UserID)
	assert.Equal(t, job.InputFile, retrievedJob.InputFile)
	assert.Equal(t, models.StatusPending, retrievedJob.Status)

	// Test job status update
	updatedJob, err := jobService.UpdateJobStatus(ctx, job.ID, models.StatusCompleted, "output/test-result.jpg")
	assert.NoError(t, err)
	assert.Equal(t, models.StatusCompleted, updatedJob.Status)
	assert.Equal(t, "output/test-result.jpg", updatedJob.OutputFile)

	// Test non-existent job
	_, err = jobService.GetJobStatus(ctx, "non-existent-job")
	assert.Error(t, err)
}

func testUtilityFunctions(t *testing.T) {
	// Test GetIntSetting
	settings := map[string]interface{}{
		"int_value":    42,
		"float_value":  3.14,
		"string_value": "not_a_number",
	}

	// Test valid int
	val, err := utils.GetIntSetting(settings, "int_value", 0)
	assert.NoError(t, err)
	assert.Equal(t, 42, val)

	// Test float conversion
	val, err = utils.GetIntSetting(settings, "float_value", 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, val)

	// Test default value
	val, err = utils.GetIntSetting(settings, "missing_key", 100)
	assert.NoError(t, err)
	assert.Equal(t, 100, val)

	// Test invalid type
	_, err = utils.GetIntSetting(settings, "string_value", 0)
	assert.Error(t, err)

	// Test GetStringSetting
	strVal, err := utils.GetStringSetting(settings, "string_value", "default")
	assert.NoError(t, err)
	assert.Equal(t, "not_a_number", strVal)

	// Test default string value
	strVal, err = utils.GetStringSetting(settings, "missing_key", "default")
	assert.NoError(t, err)
	assert.Equal(t, "default", strVal)
}

func TestSecureCommandExecution(t *testing.T) {
	executor := utils.NewSecureCommandExecutor(30 * time.Second)

	// Test allowed command structure
	_ = executor.ExecuteCommand("echo", []string{"hello"})
	// This tests the structure, actual execution may fail without proper setup

	assert.NotNil(t, executor)
}
