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

type ArchiveProcessor struct {
	config *config.Config
}

func NewArchiveProcessor(cfg *config.Config) *ArchiveProcessor {
	return &ArchiveProcessor{
		config: cfg,
	}
}

func (p *ArchiveProcessor) ProcessArchive(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Input file is already downloaded from S3 to temp by processor_s3
	inputFile := job.InputPath

	job.Progress = 30

	// Execute archive conversion
	outputFile, err := p.executeArchiveConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("archive conversion failed: %w", err)
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

func (p *ArchiveProcessor) executeArchiveConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetArchiveExtension(job.TargetFormat)
	if err != nil {
		return "", err
	}
	outputFile := filepath.Join(p.config.OutputDir,
		job.JobID+"_output"+ext)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)

	switch conversionType {
	case "RAR_TO_ZIP":
		return p.convertRarToZip(inputFile, outputFile, job)
	case "TAR_GZ_TO_ZIP":
		return p.convertTarGzToZip(inputFile, outputFile, job)
	case "ZIP_TO_TAR_GZ":
		return p.convertZipToTarGz(inputFile, outputFile, job)
	case "ZIP_TO_ZIP":
		return p.convertZipToZip(inputFile, outputFile, job)
	default:
		return "", fmt.Errorf("unsupported archive conversion: %s", conversionType)
	}
}

func (p *ArchiveProcessor) GetCompressionLevel(settings map[string]interface{}) int {
	return p.getCompressionLevel(settings)
}

func (p *ArchiveProcessor) getCompressionLevel(settings map[string]interface{}) int {
	level, err := utils.GetStringSetting(settings, "compression_level", "normal")
	if err != nil {
		return 5 // Default to normal if error occurs
	}

	levels := map[string]int{
		"store":   0, // No compression, fast
		"fast":    1, // Low compression, quick
		"normal":  5, // Balanced compression/speed
		"maximum": 7, // Best compression, slower
		"ultra":   9, // Highest compression, slowest
	}

	if lvl, exists := levels[level]; exists {
		return lvl
	}

	return 5 // Default to normal
}
