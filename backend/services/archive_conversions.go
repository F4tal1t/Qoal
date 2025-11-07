package services

import (
	"fmt"
	"os"

	"github.com/mholt/archiver/v3"
	"github.com/qoal/file-processor/models"
)

func (p *ArchiveProcessor) convertRarToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "archive_conv_")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract RAR file
	rar := archiver.NewRar()
	if err := rar.Unarchive(input, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract rar file: %w", err)
	}

	// Create ZIP file
	zip := archiver.NewZip()
	if err := zip.Archive([]string{tempDir}, output); err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertTarGzToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "archive_conv_")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract tar.gz file
	tarGz := archiver.NewTarGz()
	if err := tarGz.Unarchive(input, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract tar.gz file: %w", err)
	}

	// Create ZIP file
	zip := archiver.NewZip()
	if err := zip.Archive([]string{tempDir}, output); err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertZipToTarGz(input, output string, job *models.ProcessingJob) (string, error) {
	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "archive_conv_")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract ZIP file
	zip := archiver.NewZip()
	if err := zip.Unarchive(input, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract zip file: %w", err)
	}

	// Create tar.gz file
	tarGz := archiver.NewTarGz()
	if err := tarGz.Archive([]string{tempDir}, output); err != nil {
		return "", fmt.Errorf("failed to create tar.gz file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertZipToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Simple copy for same format
	inputData, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read input archive file: %w", err)
	}

	if err := os.WriteFile(output, inputData, 0644); err != nil {
		return "", fmt.Errorf("failed to write output archive file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) createZipFromDirectory(sourceDir, zipPath string) error {
	// Use archiver library to create ZIP from directory
	zip := archiver.NewZip()
	return zip.Archive([]string{sourceDir}, zipPath)
}
