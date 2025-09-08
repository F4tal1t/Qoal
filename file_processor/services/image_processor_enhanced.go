package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
)

type EnhancedImageProcessor struct {
	config    *config.Config
	s3Service *S3Service
}

func NewEnhancedImageProcessor(cfg *config.Config, s3 *S3Service) *EnhancedImageProcessor {
	return &EnhancedImageProcessor{
		config:    cfg,
		s3Service: s3,
	}
}

func (p *EnhancedImageProcessor) ProcessImage(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file from S3
	ext, err := utils.GetImageExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get image extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	if err := p.s3Service.DownloadFile(job.InputPath, inputFile); err != nil {
		return fmt.Errorf("failed to download input file: %w", err)
	}

	job.Progress = 30

	// Process based on conversion type
	outputFile, err := p.executeImageConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("image conversion failed: %w", err)
	}

	job.Progress = 80

	// Upload result to S3
	outputKey := fmt.Sprintf("outputs/%d/%s/converted_%s", job.UserID, job.JobID,
		filepath.Base(outputFile))
	if err := p.s3Service.UploadFile(outputFile, outputKey); err != nil {
		return fmt.Errorf("failed to upload result: %w", err)
	}

	job.OutputPath = outputKey
	job.Status = "completed"
	job.Progress = 100

	// Cleanup
	os.Remove(inputFile)
	os.Remove(outputFile)

	return nil
}

func (p *EnhancedImageProcessor) executeImageConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "JPEG_TO_PNG":
		return "", fmt.Errorf("JPEG to PNG conversion not implemented")
	case "PNG_TO_JPEG":
		return "", fmt.Errorf("PNG to JPEG conversion not implemented")
	case "PNG_TO_WEBP":
		return "", fmt.Errorf("PNG to WebP conversion not implemented")
	case "JPEG_TO_WEBP":
		return "", fmt.Errorf("JPEG to WebP conversion not implemented")
	case "WEBP_TO_JPEG":
		return "", fmt.Errorf("WebP to JPEG conversion not implemented")
	case "WEBP_TO_PNG":
		return "", fmt.Errorf("WebP to PNG conversion not implemented")
	case "HEIC_TO_JPEG":
		return "", fmt.Errorf("HEIC to JPEG conversion not implemented")
	case "BMP_TO_JPEG":
		return "", fmt.Errorf("BMP to JPEG conversion not implemented")
	case "TIFF_TO_PNG":
		return "", fmt.Errorf("TIFF to PNG conversion not implemented")
	default:
		return "", fmt.Errorf("unsupported image conversion: %s", conversionType)
	}
}
