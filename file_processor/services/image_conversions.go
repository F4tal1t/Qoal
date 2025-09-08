package services

import (
	"fmt"
	"os/exec"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
)

func (p *EnhancedImageProcessor) convertJPEGtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	args := []string{
		input,
		"-format", "png",
		"-quality", "100",
	}
	
	compressionLevel, _ := utils.GetIntSetting(job.Settings, "compression_level", 9)
	args = append(args, "-define", fmt.Sprintf("png:compression-level=%d", compressionLevel))
	args = append(args, output)
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("JPEG to PNG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertPNGtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "jpeg",
		"-quality", fmt.Sprintf("%d", quality),
		"-background", "white",
		"-flatten",
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PNG to JPEG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertPNGtoWebP(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "webp",
		"-quality", fmt.Sprintf("%d", quality),
		"-define", "webp:method=6",
		"-define", "webp:alpha-quality=100",
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PNG to WebP conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertJPEGtoWebP(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "webp",
		"-quality", fmt.Sprintf("%d", quality),
		"-define", "webp:method=6",
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("JPEG to WebP conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertWebPtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "jpeg",
		"-quality", fmt.Sprintf("%d", quality),
		"-background", "white",
		"-flatten",
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("WebP to JPEG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertWebPtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	compressionLevel, _ := utils.GetIntSetting(job.Settings, "compression_level", 9)
	
	args := []string{
		input,
		"-format", "png",
		"-define", fmt.Sprintf("png:compression-level=%d", compressionLevel),
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("WebP to PNG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertHEICtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "jpeg",
		"-quality", fmt.Sprintf("%d", quality),
		"-colorspace", "sRGB",
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("HEIC to JPEG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertBMPtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-format", "jpeg",
		"-quality", fmt.Sprintf("%d", quality),
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("BMP to JPEG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) convertTIFFtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	compressionLevel, _ := utils.GetIntSetting(job.Settings, "compression_level", 9)
	
	args := []string{
		input,
		"-format", "png",
		"-define", fmt.Sprintf("png:compression-level=%d", compressionLevel),
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("TIFF to PNG conversion failed: %w", err)
	}
	return output, nil
}

func (p *EnhancedImageProcessor) genericImageConversion(input, output string, job *models.ProcessingJob) (string, error) {
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	
	args := []string{
		input,
		"-quality", fmt.Sprintf("%d", quality),
		output,
	}
	
	cmd := exec.Command(p.config.ImageMagickPath, args...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("generic image conversion failed: %w", err)
	}
	return output, nil
}