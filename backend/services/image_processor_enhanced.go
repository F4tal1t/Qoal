package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/storage"
	"github.com/qoal/file-processor/utils"
)

type EnhancedImageProcessor struct {
	config       *config.Config
	localStorage *storage.LocalStorage
}

func NewEnhancedImageProcessor(cfg *config.Config, storage *storage.LocalStorage) *EnhancedImageProcessor {
	return &EnhancedImageProcessor{
		config:       cfg,
		localStorage: storage,
	}
}

func (p *EnhancedImageProcessor) ProcessImage(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file from storage
	ext, err := utils.GetImageExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get image extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	inputFileObj, err := p.localStorage.GetFile(job.InputPath)
	if err != nil {
		return fmt.Errorf("failed to get input file: %w", err)
	}
	defer inputFileObj.Close()

	// Copy file to temp location
	out, err := os.Create(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, inputFileObj); err != nil {
		return fmt.Errorf("failed to copy image file: %w", err)
	}

	job.Progress = 30

	// Process based on conversion type
	outputFile, err := p.executeImageConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("image conversion failed: %w", err)
	}

	job.Progress = 80

	// Save result to storage
	// Create output filename with job ID and target format
	outputFilename := fmt.Sprintf("converted_%s.%s", job.JobID, job.TargetFormat)

	// Open output file for reading
	outFile, err := os.Open(outputFile)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer outFile.Close()

	// Get file info for size
	fileInfo, err := outFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Save to processed directory using the local storage
	outputPath, err := p.localStorage.SaveFile(outFile, outputFilename, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}

	job.OutputPath = outputPath
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
