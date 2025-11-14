package config

import (
	"os"
	"strings"
)

type Config struct {
	TempDir      string
	OutputDir    string
	DatabaseURL  string
	JWTSecret    string
	RedisURL     string
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	S3Bucket     string
}

func Load() *Config {
	return &Config{
		TempDir:      os.Getenv("TEMP_DIR"),
		OutputDir:    os.Getenv("OUTPUT_DIR"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		RedisURL:     parseRedisURL(os.Getenv("REDIS_URL")),
		AWSRegion:    os.Getenv("AWS_REGION"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:     os.Getenv("AWS_S3_BUCKET"),
	}
}

// parseRedisURL converts redis:// URL to host:port format
func parseRedisURL(redisURL string) string {
	if redisURL == "" {
		return ""
	}
	// If already in host:port format, return as-is
	if !strings.HasPrefix(redisURL, "redis://") {
		return redisURL
	}
	// Remove redis:// prefix
	redisURL = strings.TrimPrefix(redisURL, "redis://")
	// Extract host:port (remove credentials if present)
	if idx := strings.Index(redisURL, "@"); idx != -1 {
		return redisURL[idx+1:]
	}
	return redisURL
}
