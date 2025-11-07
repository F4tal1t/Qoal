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

type EnhancedDocumentProcessor struct {
	config       *config.Config
	localStorage *storage.LocalStorage
}

func NewEnhancedDocumentProcessor(cfg *config.Config, storage *storage.LocalStorage) *EnhancedDocumentProcessor {
	return &EnhancedDocumentProcessor{
		config:       cfg,
		localStorage: storage,
	}
}

func (p *EnhancedDocumentProcessor) ProcessDocument(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file from storage
	ext, err := utils.GetDocumentExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get document extension: %w", err)
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
		return fmt.Errorf("failed to copy document file: %w", err)
	}

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
