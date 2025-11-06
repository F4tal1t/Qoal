package config

import "os"

type Config struct {
	TempDir     string
	OutputDir   string
	DatabaseURL string
	JWTSecret   string
	RedisURL    string
}

func Load() *Config {
	return &Config{
		TempDir:     os.Getenv("TEMP_DIR"),
		OutputDir:   os.Getenv("OUTPUT_DIR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}
