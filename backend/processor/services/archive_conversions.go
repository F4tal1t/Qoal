package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qoal/file-processor/models"
)

func (p *ArchiveProcessor) convert7ZtoZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Extract 7Z to temp directory
	tempDir := filepath.Join(p.config.TempDir, job.JobID+"_extract")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Extract using 7z command
	extractCmd := exec.Command("7z", "x", input, "-o"+tempDir, "-y")
	if err := extractCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract 7Z: %w", err)
	}

	// Create ZIP archive
	if err := p.createZipFromDirectory(tempDir, output); err != nil {
		return "", fmt.Errorf("failed to create ZIP: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertRarToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Extract RAR to temp directory
	tempDir := filepath.Join(p.config.TempDir, job.JobID+"_extract")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Extract using unrar command
	extractCmd := exec.Command("unrar", "x", input, tempDir+string(filepath.Separator))
	if err := extractCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract RAR: %w", err)
	}

	// Create ZIP archive
	if err := p.createZipFromDirectory(tempDir, output); err != nil {
		return "", fmt.Errorf("failed to create ZIP: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) convertTarGzToZip(input, output string, job *models.ProcessingJob) (string, error) {
	// Extract TAR.GZ to temp directory
	tempDir := filepath.Join(p.config.TempDir, job.JobID+"_extract")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	// Extract using tar command
	extractCmd := exec.Command("tar", "-xzf", input, "-C", tempDir)
	if err := extractCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract TAR.GZ: %w", err)
	}

	// Create ZIP archive
	if err := p.createZipFromDirectory(tempDir, output); err != nil {
		return "", fmt.Errorf("failed to create ZIP: %w", err)
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