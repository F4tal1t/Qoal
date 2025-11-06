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
	"github.com/qoal/file-processor/storage"
	"github.com/qoal/file-processor/utils"
)

type ArchiveProcessor struct {
	config       *config.Config
	localStorage *storage.LocalStorage
}

func NewArchiveProcessor(cfg *config.Config, storage *storage.LocalStorage) *ArchiveProcessor {
	return &ArchiveProcessor{
		config:       cfg,
		localStorage: storage,
	}
}

func (p *ArchiveProcessor) ProcessArchive(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file
	ext, err := utils.GetArchiveExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get archive extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	inputFileObj, err := p.localStorage.GetFile(job.InputPath)
	if err != nil {
		return fmt.Errorf("failed to get archive: %w", err)
	}
	defer inputFileObj.Close()

	// Copy file to temp location
	out, err := os.Create(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, inputFileObj); err != nil {
		return fmt.Errorf("failed to copy archive: %w", err)
	}

	job.Progress = 30

	// Execute archive conversion
	outputFile, err := p.executeArchiveConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("archive conversion failed: %w", err)
	}

	job.Progress = 80

	// Save result to storage
	outputKey := fmt.Sprintf("outputs/%s/%s/converted_%s", job.UserID, job.JobID,
		filepath.Base(outputFile))

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

	if _, err := p.localStorage.SaveFile(outFile, outputKey, fileInfo.Size()); err != nil {
		return fmt.Errorf("failed to save archive result: %w", err)
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
		return p.convert7ZtoZip(inputFile, outputFile, job)
	case "RAR_TO_ZIP":
		return p.convertRarToZip(inputFile, outputFile, job)
	case "TAR_GZ_TO_ZIP":
		return p.convertTarGzToZip(inputFile, outputFile, job)
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
