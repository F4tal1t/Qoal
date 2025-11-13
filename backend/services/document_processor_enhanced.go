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

type EnhancedDocumentProcessor struct {
	config *config.Config
}

func NewEnhancedDocumentProcessor(cfg *config.Config) *EnhancedDocumentProcessor {
	return &EnhancedDocumentProcessor{
		config: cfg,
	}
}

func (p *EnhancedDocumentProcessor) ProcessDocument(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Input file is already downloaded from S3 to temp by processor_s3
	inputFile := job.InputPath

	job.Progress = 30

	// Validate document
	if err := p.validateDocument(inputFile); err != nil {
		return fmt.Errorf("document validation failed: %w", err)
	}

	job.Progress = 40

	// Process based on conversion type
	outputFile, err := p.executeDocumentConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("document conversion failed: %w", err)
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

func (p *EnhancedDocumentProcessor) executeDocumentConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetDocumentExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "TEXT_TO_PDF":
		return p.convertTextToPDF(inputFile, outputFile, job)
	case "DOCX_TO_TEXT":
		return p.convertDocxToText(inputFile, outputFile, job)
	case "TEXT_TO_DOCX":
		return p.convertTextToDocx(inputFile, outputFile, job)
	case "XLSX_TO_CSV":
		return p.convertXlsxToCSV(inputFile, outputFile, job)
	case "CSV_TO_XLSX":
		return p.convertCSVToXlsx(inputFile, outputFile, job)
	default:
		return "", fmt.Errorf("unsupported document conversion: %s", conversionType)
	}
}

func (p *EnhancedDocumentProcessor) validateDocument(filePath string) error {
	// Check file size limits
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fileInfo.Size() > 50*1024*1024 { // 50MB limit
		return fmt.Errorf("document too large: %d bytes", fileInfo.Size())
	}

	// Basic file validation
	return nil
}
