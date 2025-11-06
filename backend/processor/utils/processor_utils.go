package utils

import (
	"fmt"
	"strings"
)

func GetArchiveExtension(format string) (string, error) {
	switch strings.ToLower(format) {
	case "zip":
		return ".zip", nil
	case "7z":
		return ".7z", nil
	case "rar":
		return ".rar", nil
	case "tar.gz":
		return ".tar.gz", nil
	default:
		return "", fmt.Errorf("unsupported archive format: %s", format)
	}
}

func GetAudioExtension(format string) (string, error) {
	switch strings.ToLower(format) {
	case "mp3":
		return ".mp3", nil
	case "wav":
		return ".wav", nil
	case "flac":
		return ".flac", nil
	case "aac":
		return ".aac", nil
	case "m4a":
		return ".m4a", nil
	case "ogg":
		return ".ogg", nil
	default:
		return "", fmt.Errorf("unsupported audio format: %s", format)
	}
}

func GetStringSetting(settings map[string]interface{}, key, defaultValue string) (string, error) {
	if val, exists := settings[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal, nil
		}
		return "", fmt.Errorf("invalid type for setting %s", key)
	}
	return defaultValue, nil
}

func GetIntSetting(settings map[string]interface{}, key string, defaultValue int) (int, error) {
	if val, exists := settings[key]; exists {
		switch v := val.(type) {
		case int:
			return v, nil
		case float64:
			return int(v), nil
		default:
			return 0, fmt.Errorf("invalid type for setting %s", key)
		}
	}
	return defaultValue, nil
}

func GetImageExtension(format string) (string, error) {
	switch strings.ToLower(format) {
	case "jpeg":
		return ".jpg", nil
	case "jpg":
		return ".jpg", nil
	case "png":
		return ".png", nil
	case "webp":
		return ".webp", nil
	case "heic":
		return ".heic", nil
	case "bmp":
		return ".bmp", nil
	case "tiff":
		return ".tiff", nil
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}
}

func GetDocumentExtension(format string) (string, error) {
	switch strings.ToLower(format) {
	case "pdf":
		return ".pdf", nil
	case "docx":
		return ".docx", nil
	case "doc":
		return ".doc", nil
	case "txt":
		return ".txt", nil
	case "rtf":
		return ".rtf", nil
	case "odt":
		return ".odt", nil
	default:
		return "", fmt.Errorf("unsupported document format: %s", format)
	}
}

func GetVideoExtension(format string) (string, error) {
	switch strings.ToLower(format) {
	case "mp4":
		return ".mp4", nil
	case "avi":
		return ".avi", nil
	case "mov":
		return ".mov", nil
	case "mkv":
		return ".mkv", nil
	case "webm":
		return ".webm", nil
	case "flv":
		return ".flv", nil
	case "wmv":
		return ".wmv", nil
	default:
		return "", fmt.Errorf("unsupported video format: %s", format)
	}
}