package services

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// decodeImage decodes an image file supporting multiple formats
func (p *EnhancedImageProcessor) decodeImage(inputFile string) (image.Image, string, error) {
	input, err := os.Open(inputFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	// Determine file format from extension
	ext := strings.ToLower(filepath.Ext(inputFile))

	switch ext {
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(input)
		return img, "jpeg", err
	case ".png":
		img, err := png.Decode(input)
		return img, "png", err
	case ".gif":
		img, err := gif.Decode(input)
		return img, "gif", err
	case ".bmp":
		img, err := bmp.Decode(input)
		return img, "bmp", err
	case ".tiff", ".tif":
		img, err := tiff.Decode(input)
		return img, "tiff", err
	default:
		// Try to decode with default image.Decode
		img, format, err := image.Decode(input)
		return img, format, err
	}
}

func (p *EnhancedImageProcessor) convertJPEGtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode image (handles JPEG and other formats)
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Get compression level from settings
	compressionLevel, _ := utils.GetIntSetting(job.Settings, "compression_level", 9)
	compressionLevel = 9 - compressionLevel // Invert for PNG compression (0=fast, 9=best)
	if compressionLevel < 0 || compressionLevel > 9 {
		compressionLevel = 0
	}

	// Encode as PNG
	encoder := png.Encoder{
		CompressionLevel: png.CompressionLevel(compressionLevel),
	}
	if err := encoder.Encode(out, img); err != nil {
		return "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertPNGtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode image (handles PNG and other formats)
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Get quality from settings
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	if quality < 1 || quality > 100 {
		quality = 85
	}

	// Encode as JPEG
	options := jpeg.Options{Quality: quality}
	if err := jpeg.Encode(out, img, &options); err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertBMPtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode BMP image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode BMP: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Get quality from settings
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	if quality < 1 || quality > 100 {
		quality = 85
	}

	// Encode as JPEG
	options := jpeg.Options{Quality: quality}
	if err := jpeg.Encode(out, img, &options); err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertTIFFtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode TIFF image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode TIFF: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as PNG with compression
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err := encoder.Encode(out, img); err != nil {
		return "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertJPEGtoBMP(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode JPEG image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as BMP
	if err := bmp.Encode(out, img); err != nil {
		return "", fmt.Errorf("failed to encode BMP: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertPNGtoBMP(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode PNG image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as BMP
	if err := bmp.Encode(out, img); err != nil {
		return "", fmt.Errorf("failed to encode BMP: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertJPEGtoTIFF(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode JPEG image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode JPEG: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as TIFF
	if err := tiff.Encode(out, img, nil); err != nil {
		return "", fmt.Errorf("failed to encode TIFF: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertPNGtoTIFF(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode PNG image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as TIFF
	if err := tiff.Encode(out, img, nil); err != nil {
		return "", fmt.Errorf("failed to encode TIFF: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertGIFtoJPEG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode GIF image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode GIF: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Get quality from settings
	quality, _ := utils.GetIntSetting(job.Settings, "quality", 85)
	if quality < 1 || quality > 100 {
		quality = 85
	}

	// Encode as JPEG (first frame only for GIF)
	options := jpeg.Options{Quality: quality}
	if err := jpeg.Encode(out, img, &options); err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) convertGIFtoPNG(input, output string, job *models.ProcessingJob) (string, error) {
	// Decode GIF image
	img, _, err := p.decodeImage(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode GIF: %w", err)
	}

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Encode as PNG
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err := encoder.Encode(out, img); err != nil {
		return "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	return output, nil
}

func (p *EnhancedImageProcessor) genericImageConversion(inputFile, outputFile string, job *models.ProcessingJob) (string, error) {
	// Decode image using our universal decoder
	img, format, err := p.decodeImage(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Apply processing options if specified
	if job.Settings != nil {
		// Resize if specified
		if width, ok := job.Settings["width"].(float64); ok {
			if height, ok := job.Settings["height"].(float64); ok {
				img = imaging.Resize(img, int(width), int(height), imaging.Lanczos)
			}
		}

		// Crop if specified
		if cropWidth, ok := job.Settings["crop_width"].(float64); ok {
			if cropHeight, ok := job.Settings["crop_height"].(float64); ok {
				if cropX, ok := job.Settings["crop_x"].(float64); ok {
					if cropY, ok := job.Settings["crop_y"].(float64); ok {
						img = imaging.Crop(img, image.Rect(int(cropX), int(cropY), int(cropX+cropWidth), int(cropY+cropHeight)))
					}
				}
			}
		}

		// Rotate if specified
		if rotate, ok := job.Settings["rotate"].(float64); ok {
			switch rotate {
			case 90:
				img = imaging.Rotate90(img)
			case 180:
				img = imaging.Rotate180(img)
			case 270:
				img = imaging.Rotate270(img)
			}
		}
	}

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	// Determine output format from extension
	outputExt := strings.ToLower(filepath.Ext(outputFile))

	switch outputExt {
	case ".jpg", ".jpeg":
		quality := 90
		if job.Settings != nil {
			if q, ok := job.Settings["quality"].(float64); ok {
				quality = int(q)
			}
		}
		if quality < 1 || quality > 100 {
			quality = 90
		}
		options := &jpeg.Options{Quality: quality}
		err = jpeg.Encode(output, img, options)
	case ".png":
		encoder := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		err = encoder.Encode(output, img)
	case ".gif":
		err = gif.Encode(output, img, nil)
	case ".bmp":
		err = bmp.Encode(output, img)
	case ".tiff", ".tif":
		err = tiff.Encode(output, img, nil)
	default:
		// Try to use the original format if possible
		switch format {
		case "jpeg":
			quality := 90
			if job.Settings != nil {
				if q, ok := job.Settings["quality"].(float64); ok {
					quality = int(q)
				}
			}
			options := &jpeg.Options{Quality: quality}
			err = jpeg.Encode(output, img, options)
		case "png":
			encoder := png.Encoder{
				CompressionLevel: png.BestCompression,
			}
			err = encoder.Encode(output, img)
		default:
			return "", fmt.Errorf("unsupported output format: %s", outputExt)
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return outputFile, nil
}
