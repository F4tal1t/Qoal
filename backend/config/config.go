package config

import "os"

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
		RedisURL:     os.Getenv("REDIS_URL"),
		AWSRegion:    os.Getenv("AWS_REGION"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:     os.Getenv("AWS_S3_BUCKET"),
	}
}
