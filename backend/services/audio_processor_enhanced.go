package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
)

type EnhancedAudioProcessor struct {
	config *config.Config
}

func NewEnhancedAudioProcessor(cfg *config.Config) *EnhancedAudioProcessor {
	return &EnhancedAudioProcessor{
		config: cfg,
	}
}

func (p *EnhancedAudioProcessor) ProcessAudio(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Input file is already downloaded from S3 to temp by processor_s3
	inputFile := job.InputPath

	job.Progress = 30

	// Execute audio conversion
	outputFile, err := p.executeAudioConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("audio conversion failed: %w", err)
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
