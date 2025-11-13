package services

import (
	"fmt"
	"os"
	"path/filepath"

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
	config *config.Config
}

func NewEnhancedVideoProcessor(cfg *config.Config) *EnhancedVideoProcessor {
	return &EnhancedVideoProcessor{
		config: cfg,
	}
}

func (p *EnhancedVideoProcessor) ProcessVideo(job *models.ProcessingJob) error {
	job.Status = "processing"
	job.Progress = 10

	// Input file is already downloaded from S3 to temp by processor_s3
	inputFile := job.InputPath

	job.Progress = 30

	// Process based on conversion type
	outputFile, err := p.executeVideoConversion(inputFile, job)
	if err != nil {
		return fmt.Errorf("video conversion failed: %w", err)
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

func (p *EnhancedVideoProcessor) executeVideoConversion(inputFile string, job *models.ProcessingJob) (string, error) {
	ext, err := utils.GetVideoExtension(job.TargetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to get target extension: %w", err)
	}
	outputFile := filepath.Join(p.config.OutputDir, job.JobID+"_output"+ext)

	conversionType := job.SourceFormat + "_TO_" + job.TargetFormat
	switch conversionType {
	case "MP4_TO_AVI":
		return p.convertMP4toAVI(inputFile, outputFile, job)
	case "AVI_TO_MP4":
		return p.convertAVItoMP4(inputFile, outputFile, job)
	case "MP4_TO_MOV":
		return p.convertMP4toMOV(inputFile, outputFile, job)
	case "MOV_TO_MP4":
		return p.convertMOVtoMP4(inputFile, outputFile, job)
	case "MP4_TO_WEBM":
		return p.convertMP4toWEBM(inputFile, outputFile, job)
	case "WEBM_TO_MP4":
		return p.convertWEBMtoMP4(inputFile, outputFile, job)
	case "MP4_TO_MKV":
		return p.convertMP4toMKV(inputFile, outputFile, job)
	case "MKV_TO_MP4":
		return p.convertMKVtoMP4(inputFile, outputFile, job)
	default:
		return "", fmt.Errorf("unsupported video conversion: %s", conversionType)
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
