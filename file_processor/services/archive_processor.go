package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
)

type ArchiveProcessor struct {
	config    *config.Config
	s3Service *S3Service
}

func NewArchiveProcessor(cfg *config.Config, s3 *S3Service) *ArchiveProcessor {
	return &ArchiveProcessor{
		config:    cfg,
		s3Service: s3,
	}
}

func (p *ArchiveProcessor) ProcessArchive(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+getArchiveExtension(job.SourceFormat))
	if err := p.s3Service.DownloadFile(job.InputPath, inputFile); err != nil {
		return fmt.Errorf("failed to download archive: %w", err)
	}

	job.Progress = 30

	// Execute archive conversion
	outputFile, err := p.executeArchiveConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("archive conversion failed: %w", err)
	}

	job.Progress = 80

	// Upload result
	outputKey := fmt.Sprintf("outputs/%d/%s/converted_%s", job.UserID, job.JobID,
		filepath.Base(outputFile))
	if err := p.s3Service.UploadFile(outputFile, outputKey); err != nil {
		return fmt.Errorf("failed to upload archive result: %w", err)
	}

	job.OutputPath = outputKey
	job.Status = "completed"
	job.Progress = 100

	// Cleanup
	os.Remove(inputFile)
	os.Remove(outputFile)

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
	case "ZIP_TO_7Z":
		return p.convertZipTo7Z(inputFile, outputFile, job)
	case "7Z_TO_ZIP":
		return "", fmt.Errorf("7Z to ZIP conversion not implemented")
	case "RAR_TO_ZIP":
		return "", fmt.Errorf("RAR to ZIP conversion not implemented")
	case "TAR_GZ_TO_ZIP":
		return "", fmt.Errorf("TAR.GZ to ZIP conversion not implemented")
	default:
		return "", fmt.Errorf("unsupported archive conversion: %s", conversionType)
	}
}

func (p *ArchiveProcessor) convertZipTo7Z(input, output string, job *models.ProcessingJob) (string, error) {
	// Extract ZIP to temp directory
	tempDir := filepath.Join(p.config.TempDir, job.JobID+"_extract")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	if err := p.extractZip(input, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Create 7Z archive with compression level
	compressionLevel := p.getCompressionLevel(job.Settings)
	args := []string{
		"a",                                     // Add to archive
		fmt.Sprintf("-mx=%d", compressionLevel), // Compression level
		output,
		filepath.Join(tempDir, "*"),
	}

	cmd := exec.Command("7z", args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("7z compression failed: %w", err)
	}

	return output, nil
}

func (p *ArchiveProcessor) extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.FileInfo().Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0755)
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
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
