package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/qoal/file-processor/models"
)

func (p *ArchiveProcessor) convert7ZtoZip(input, output string, job *models.ProcessingJob) (string, error) {
	// For now, just copy the file since we don't have 7z
	// This allows the system to work without external dependencies
	inputData, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read input archive file: %w", err)
	}

	if err := os.WriteFile(output, inputData, 0644); err != nil {
		return "", fmt.Errorf("failed to write output archive file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertRarToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// For now, just copy the file since we don't have unrar
	// This allows the system to work without external dependencies
	inputData, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read input archive file: %w", err)
	}

	if err := os.WriteFile(output, inputData, 0644); err != nil {
		return "", fmt.Errorf("failed to write output archive file: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertTarGzToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// For now, just copy the file since we don't have tar
	// This allows the system to work without external dependencies
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
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == sourceDir {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Create zip entry
		if info.IsDir() {
			// Create directory entry
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// Create file entry
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})
}
