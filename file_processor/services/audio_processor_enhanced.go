package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/models"
)

type EnhancedAudioProcessor struct {
	config    *config.Config
	s3Service *S3Service
}

func NewEnhancedAudioProcessor(cfg *config.Config, s3 *S3Service) *EnhancedAudioProcessor {
	return &EnhancedAudioProcessor{
		config:    cfg,
		s3Service: s3,
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
	if err := p.s3Service.DownloadFile(job.InputPath, inputFile); err != nil {
		return fmt.Errorf("failed to download audio file: %w", err)
	}

	job.Progress = 30

	// Execute audio conversion
	outputFile, err := p.executeAudioConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("audio conversion failed: %w", err)
	}

	job.Progress = 80

	// Upload result
	outputKey := fmt.Sprintf("outputs/%d/%s/converted_%s", job.UserID, job.JobID,
		filepath.Base(outputFile))
	if err := p.s3Service.UploadFile(outputFile, outputKey); err != nil {
		return fmt.Errorf("failed to upload audio result: %w", err)
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

	// Build conversion command based on specific conversion type
	cmd, err := p.buildSpecificAudioCommand(inputFile, outputFile, job)
	if err != nil {
		return "", err
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("audio conversion command failed: %w", err)
	}

	return outputFile, nil
}

func (p *EnhancedAudioProcessor) buildSpecificAudioCommand(input, output string, job *models.ProcessingJob) (*exec.Cmd, error) {
	args := []string{"-i", input, "-y"}

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)

	switch conversionType {
	case "MP3_TO_WAV":
		// Lossy to lossless conversion
		args = append(args, "-c:a", "pcm_s16le", "-ar", "44100")

	case "WAV_TO_MP3":
		// Lossless to lossy with quality control
		quality := p.getAudioQualityPreset(job.Settings)
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality.Bitrate)

	case "FLAC_TO_MP3":
		// High quality lossless to portable lossy
		quality := p.getAudioQualityPreset(job.Settings)
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality.Bitrate)

	case "WAV_TO_FLAC":
		// Lossless compression without quality loss
		args = append(args, "-c:a", "flac", "-compression_level", "8")

	case "AAC_TO_MP3":
		// Cross-platform compatibility conversion
		quality := p.getAudioQualityPreset(job.Settings)
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality.Bitrate)

	case "M4A_TO_MP3":
		// Apple to universal format
		quality := p.getAudioQualityPreset(job.Settings)
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality.Bitrate)

	case "OGG_TO_MP3":
		// Open source to standard format
		quality := p.getAudioQualityPreset(job.Settings)
		args = append(args, "-c:a", "libmp3lame", "-b:a", quality.Bitrate)

	default:
		return nil, fmt.Errorf("unsupported audio conversion: %s", conversionType)
	}

	args = append(args, output)
	return exec.Command(p.config.FFmpegPath, args...), nil
}

type AudioQuality struct {
	Bitrate     string
	SampleRate  int
	Description string
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
