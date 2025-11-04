package config

import "os"

type Config struct {
	ImageMagickPath string
	FFmpegPath      string
	TempDir         string
	OutputDir       string
}

func Load() *Config {
	return &Config{
		ImageMagickPath: os.Getenv("IMAGE_MAGICK_PATH"),
		FFmpegPath:      os.Getenv("FFMPEG_PATH"),
		TempDir:         os.Getenv("TEMP_DIR"),
		OutputDir:       os.Getenv("OUTPUT_DIR"),
	}
}
