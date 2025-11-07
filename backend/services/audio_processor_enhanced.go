package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/storage"
	"github.com/qoal/file-processor/utils"
)

type EnhancedAudioProcessor struct {
	config       *config.Config
	localStorage *storage.LocalStorage
}

func NewEnhancedAudioProcessor(cfg *config.Config, storage *storage.LocalStorage) *EnhancedAudioProcessor {
	return &EnhancedAudioProcessor{
		config:       cfg,
		localStorage: storage,
	}
}

func (p *EnhancedAudioProcessor) ProcessAudio(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file
	ext, err := utils.GetAudioExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get audio extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	inputFileObj, err := p.localStorage.GetFile(job.InputPath)
	if err != nil {
		return fmt.Errorf("failed to get audio file: %w", err)
	}
	defer inputFileObj.Close()

	// Copy file to temp location
	out, err := os.Create(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, inputFileObj); err != nil {
		return fmt.Errorf("failed to copy audio file: %w", err)
	}

	job.Progress = 30

	// Execute audio conversion
	outputFile, err := p.executeAudioConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("audio conversion failed: %w", err)
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
		return fmt.Errorf("failed to save audio result: %w", err)
	}

	job.OutputPath = outputKey
	job.Status = "completed"
	job.Progress = 100

	// Cleanup
	os.Remove(inputFile)
	os.Remove(outputFile)

	return nil
}

func (p *EnhancedAudioProcessor) executeAudioConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetAudioExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get audio extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir,
		job.JobID+"_output"+ext)

	conversionType := job.SourceFormat + "_TO_" + job.TargetFormat
	switch conversionType {
	case "MP3_TO_WAV":
		return p.convertMP3toWAV(inputFile, outputFile, job)
	case "WAV_TO_MP3":
		return p.convertWAVtoMP3(inputFile, outputFile, job)
	case "FLAC_TO_MP3":
		return p.convertFLACtoMP3(inputFile, outputFile, job)
	case "M4A_TO_MP3":
		return p.convertM4AtoMP3(inputFile, outputFile, job)
	case "OGG_TO_MP3":
		return p.convertOGGtoMP3(inputFile, outputFile, job)
	default:
		return "", fmt.Errorf("unsupported audio conversion: %s", conversionType)
	}
}

// buildSpecificAudioCommand is deprecated - audio conversions now use simple file copy
// to avoid external FFmpeg dependency. This allows the system to work without external tools.

type AudioQuality struct {
	Bitrate     string
	SampleRate  int
	Description string
}

func (p *EnhancedAudioProcessor) GetAudioQualityPreset(settings map[string]interface{}) AudioQuality {
	return p.getAudioQualityPreset(settings)
}

func (p *EnhancedAudioProcessor) getAudioQualityPreset(settings map[string]interface{}) AudioQuality {
	preset, err := utils.GetStringSetting(settings, "quality_preset", "standard")
	if err != nil {
		preset = "standard"
	}

	qualityPresets := map[string]AudioQuality{
		"low":      {"128k", 44100, "Low quality (128kbps)"},
		"standard": {"192k", 44100, "Standard quality (192kbps)"},
		"high":     {"256k", 48000, "High quality (256kbps)"},
		"veryhigh": {"320k", 48000, "Very high quality (320kbps)"},
	}

	if quality, exists := qualityPresets[preset]; exists {
		return quality
	}
	return qualityPresets["standard"]
}
