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

type VideoPreset struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Bitrate     string `json:"bitrate"`
	Description string `json:"description"`
}

type EnhancedVideoProcessor struct {
	config    *config.Config
	s3Service *S3Service
}

func NewEnhancedVideoProcessor(cfg *config.Config, s3 *S3Service) *EnhancedVideoProcessor {
	return &EnhancedVideoProcessor{
		config:    cfg,
		s3Service: s3,
	}
}

func (p *EnhancedVideoProcessor) ProcessVideo(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Download input file from S3
	ext, err := utils.GetVideoExtension(job.SourceFormat)
	if err != nil {
		return fmt.Errorf("failed to get video extension: %w", err)
	}
	inputFile := filepath.Join(p.config.TempDir, job.JobID+"_input"+ext)
	if err := p.s3Service.DownloadFile(job.InputPath, inputFile); err != nil {
		return fmt.Errorf("failed to download input file: %w", err)
	}

	job.Progress = 30

	// Process based on conversion type
	outputFile, err := p.executeVideoConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("video conversion failed: %w", err)
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

func (p *EnhancedVideoProcessor) executeVideoConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetVideoExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	// Get video preset
	preset := p.getVideoPreset(job.Settings)

	conversionType := strings.ToUpper(job.SourceFormat) + "_TO_" + strings.ToUpper(job.TargetFormat)
	switch conversionType {
	case "MP4_TO_AVI":
		return p.convertMP4toAVI(inputFile, outputFile, preset, job)
	case "MOV_TO_MP4":
		return p.convertMOVtoMP4(inputFile, outputFile, preset, job)
	default:
		return p.genericVideoConversion(inputFile, outputFile, preset, job)
	}
}

func (p *EnhancedVideoProcessor) getVideoPreset(settings map[string]interface{}) VideoPreset {
	presetName := "720p"
	if settings != nil {
		if preset, ok := settings["resolution_preset"].(string); ok {
			presetName = preset
		}
	}

	presets := map[string]VideoPreset{
		"4k": {
			Width: 3840, Height: 2160, Bitrate: "8000k",
			Description: "4K Ultra HD (large files, best quality)",
		},
		"1080p": {
			Width: 1920, Height: 1080, Bitrate: "2000k",
			Description: "Full HD (standard high quality)",
		},
		"720p": {
			Width: 1280, Height: 720, Bitrate: "1000k",
			Description: "HD (balanced size/quality)",
		},
		"480p": {
			Width: 854, Height: 480, Bitrate: "500k",
			Description: "SD (web optimized, smaller files)",
		},
		"360p": {
			Width: 640, Height: 360, Bitrate: "300k",
			Description: "Mobile optimized (low bandwidth)",
		},
	}

	if preset, exists := presets[presetName]; exists {
		return preset
	}
	return presets["720p"] // Default
}

func (p *EnhancedVideoProcessor) convertMP4toAVI(inputFile, outputFile string, preset VideoPreset, job *models.ProcessingJob) (string, error) {
	args := []string{
		"-i", inputFile,
		"-vf", fmt.Sprintf("scale=%d:%d", preset.Width, preset.Height),
		"-b:v", preset.Bitrate,
		"-c:v", "libx264",
		"-c:a", "mp3",
		"-y", outputFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	return outputFile, cmd.Run()
}

func (p *EnhancedVideoProcessor) convertMOVtoMP4(inputFile, outputFile string, preset VideoPreset, job *models.ProcessingJob) (string, error) {
	args := []string{
		"-i", inputFile,
		"-vf", fmt.Sprintf("scale=%d:%d", preset.Width, preset.Height),
		"-b:v", preset.Bitrate,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-y", outputFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	return outputFile, cmd.Run()
}

func (p *EnhancedVideoProcessor) genericVideoConversion(inputFile, outputFile string, preset VideoPreset, job *models.ProcessingJob) (string, error) {
	args := []string{
		"-i", inputFile,
		"-vf", fmt.Sprintf("scale=%d:%d", preset.Width, preset.Height),
		"-b:v", preset.Bitrate,
		"-y", outputFile,
	}

	cmd := exec.Command("ffmpeg", args...)
	return outputFile, cmd.Run()
}
