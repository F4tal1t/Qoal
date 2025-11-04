package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
)

type EnhancedDocumentProcessor struct {
	config    *config.Config
	s3Service *S3Service
}

func NewEnhancedDocumentProcessor(cfg *config.Config, s3 *S3Service) *EnhancedDocumentProcessor {
	return &EnhancedDocumentProcessor{
		config:    cfg,
		s3Service: s3,
	}
}

func (p *EnhancedDocumentProcessor) ProcessDocument(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file from S3
	ext, err := utils.GetDocumentExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get document extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	if err := p.s3Service.DownloadFile(job.InputPath, inputFile); err != nil {
		return fmt.Errorf("failed to download input file: %w", err)
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

func (p *EnhancedDocumentProcessor) executeDocumentConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetDocumentExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "PDF_TO_DOCX":
		return p.convertPDFToDocx(inputFile, outputFile, job)
	case "DOCX_TO_PDF":
		return p.convertDocxToPDF(inputFile, outputFile, job)
	case "XLSX_TO_CSV":
		return p.convertXlsxToCSV(inputFile, outputFile, job)
	default:
		return p.genericDocumentConversion(inputFile, outputFile, job)
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

func (p *EnhancedDocumentProcessor) convertPDFToDocx(inputFile, outputFile string, job *models.ProcessingJob) (string, error) {
	// Most requested conversion - PDF to editable Word document
	args := []string{
		"--headless",
		"--convert-to", "docx",
		"--outdir", filepath.Dir(outputFile),
		inputFile,
	}

	cmd := exec.Command("libreoffice", args...)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// LibreOffice creates file with original name + new extension
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	generatedFile := filepath.Join(filepath.Dir(outputFile), baseName+".docx")
	
	// Rename to expected output file
	if err := os.Rename(generatedFile, outputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}

func (p *EnhancedDocumentProcessor) convertDocxToPDF(inputFile, outputFile string, job *models.ProcessingJob) (string, error) {
	args := []string{
		"--headless",
		"--convert-to", "pdf",
		"--outdir", filepath.Dir(outputFile),
		inputFile,
	}

	cmd := exec.Command("libreoffice", args...)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	generatedFile := filepath.Join(filepath.Dir(outputFile), baseName+".pdf")
	
	if err := os.Rename(generatedFile, outputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}

func (p *EnhancedDocumentProcessor) convertXlsxToCSV(inputFile, outputFile string, job *models.ProcessingJob) (string, error) {
	args := []string{
		"--headless",
		"--convert-to", "csv",
		"--outdir", filepath.Dir(outputFile),
		inputFile,
	}

	cmd := exec.Command("libreoffice", args...)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	generatedFile := filepath.Join(filepath.Dir(outputFile), baseName+".csv")
	
	if err := os.Rename(generatedFile, outputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}

func (p *EnhancedDocumentProcessor) genericDocumentConversion(inputFile, outputFile string, job *models.ProcessingJob) (string, error) {
	targetExt := strings.TrimPrefix(filepath.Ext(outputFile), ".")
	
	args := []string{
		"--headless",
		"--convert-to", targetExt,
		"--outdir", filepath.Dir(outputFile),
		inputFile,
	}

	cmd := exec.Command("libreoffice", args...)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	generatedFile := filepath.Join(filepath.Dir(outputFile), baseName+"."+targetExt)
	
	if err := os.Rename(generatedFile, outputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}