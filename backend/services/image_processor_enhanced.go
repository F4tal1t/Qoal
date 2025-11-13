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
	config *config.Config
}

func NewEnhancedImageProcessor(cfg *config.Config) *EnhancedImageProcessor {
	return &EnhancedImageProcessor{
		config: cfg,
	}
}

func (p *EnhancedImageProcessor) ProcessImage(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Input file is already downloaded from S3 to temp by processor_s3
	inputFile := job.InputPath

	job.Progress = 30

	// Process based on conversion type
	outputFile, err := p.executeImageConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("image conversion failed: %w", err)
	}

	job.Progress = 80

	// Set output path (will be uploaded by S3 processor if needed)
	job.OutputPath = outputFile
	job.Status = "completed"
	job.Progress = 100

	// Cleanup input only (output will be cleaned by S3 processor)
	if inputFile != job.InputPath {
		os.Remove(inputFile)
	}

	return nil
}

func (p *EnhancedImageProcessor) executeImageConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetImageExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}

	// Ensure output directory exists
	os.MkdirAll(p.config.OutputDir, 0755)
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "JPEG_TO_PNG":
		return p.convertJPEGtoPNG(inputFile, outputFile, job)
	case "PNG_TO_JPEG":
		return p.convertPNGtoJPEG(inputFile, outputFile, job)
	case "BMP_TO_JPEG":
		return p.convertBMPtoJPEG(inputFile, outputFile, job)
	case "TIFF_TO_PNG":
		return p.convertTIFFtoPNG(inputFile, outputFile, job)
	case "JPEG_TO_BMP":
		return p.convertJPEGtoBMP(inputFile, outputFile, job)
	case "PNG_TO_BMP":
		return p.convertPNGtoBMP(inputFile, outputFile, job)
	case "JPEG_TO_TIFF":
		return p.convertJPEGtoTIFF(inputFile, outputFile, job)
	case "PNG_TO_TIFF":
		return p.convertPNGtoTIFF(inputFile, outputFile, job)
	case "GIF_TO_JPEG":
		return p.convertGIFtoJPEG(inputFile, outputFile, job)
	case "GIF_TO_PNG":
		return p.convertGIFtoPNG(inputFile, outputFile, job)
	case "PNG_TO_WEBP", "JPEG_TO_WEBP":
		// For now, fall back to generic conversion for WebP since we need additional libraries
		return p.genericImageConversion(inputFile, outputFile, job)
	case "WEBP_TO_JPEG", "WEBP_TO_PNG":
		// For now, fall back to generic conversion for WebP since we need additional libraries
		return p.genericImageConversion(inputFile, outputFile, job)
	case "HEIC_TO_JPEG":
		// These formats require special handling, use generic for now
		return p.genericImageConversion(inputFile, outputFile, job)
	default:
		return p.genericImageConversion(inputFile, outputFile, job)
	}
}
