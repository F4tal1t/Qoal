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
	ext, err := utils.GetImageExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "JPEG_TO_PNG":
		return p.convertJPEGtoPNG(inputFile, outputFile, job)
	case "PNG_TO_JPEG":
		return p.convertPNGtoJPEG(inputFile, outputFile, job)
	case "PNG_TO_WEBP":
		return p.convertPNGtoWebP(inputFile, outputFile, job)
	case "JPEG_TO_WEBP":
		return p.convertJPEGtoWebP(inputFile, outputFile, job)
	case "WEBP_TO_JPEG":
		return p.convertWebPtoJPEG(inputFile, outputFile, job)
	case "WEBP_TO_PNG":
		return p.convertWebPtoPNG(inputFile, outputFile, job)
	case "HEIC_TO_JPEG":
		return p.convertHEICtoJPEG(inputFile, outputFile, job)
	case "BMP_TO_JPEG":
		return p.convertBMPtoJPEG(inputFile, outputFile, job)
	case "TIFF_TO_PNG":
		return p.convertTIFFtoPNG(inputFile, outputFile, job)
	default:
		return p.genericImageConversion(inputFile, outputFile, job)
	}
}
