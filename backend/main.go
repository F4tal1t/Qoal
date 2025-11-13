package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/handlers"
	"github.com/qoal/file-processor/middleware"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/storage"
	"github.com/qoal/file-processor/utils"
	"github.com/qoal/file-processor/worker"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to PostgreSQL database")

	// Run database migrations
	if err := utils.MigrateDatabase(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed database with initial data
	if err := utils.SeedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis - processing jobs will be disabled")
		redisClient = nil
	}

	// Initialize services
	authService := services.NewAuthService(db)

	// Initialize S3 storage
	s3Storage, err := storage.NewS3Storage(cfg.AWSRegion, cfg.S3Bucket, cfg.AWSAccessKey, cfg.AWSSecretKey)
	if err != nil {
		log.Fatal("Failed to initialize S3 storage:", err)
	}
	log.Println("Using S3 storage")

	var jobService *services.JobService
	if redisClient != nil {
		jobService = services.NewJobService(db, redisClient)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	var jobHandler *handlers.JobHandler
	if jobService != nil {
		jobHandler = handlers.NewJobHandler(jobService)
	}

	// Initialize upload handler
	uploadHandler := handlers.NewUploadHandlerS3(db, s3Storage, jobService)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuth(authService))
	{
		protected.GET("/auth/profile", authHandler.GetProfile)
		if jobHandler != nil {
			protected.POST("/process", jobHandler.CreateJobHandler)
			protected.GET("/status/:id", jobHandler.GetJobStatusHandler)
		}
		
		protected.POST("/upload", func(c *gin.Context) {
			uploadHandler.UploadFileS3(c, s3Storage)
		})
		protected.GET("/download/:id", func(c *gin.Context) {
			uploadHandler.DownloadFileS3(c, s3Storage)
		})
		
		protected.GET("/jobs", uploadHandler.GetUserJobs)
		protected.GET("/jobs/:id", uploadHandler.GetJobStatus)
	}

	// Start S3 worker in background (only if Redis is available)
	if jobService != nil {
		go func() {
			processor := worker.NewProcessorS3(jobService, cfg, redisClient, s3Storage)
			processor.Start(ctx)
		}()
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
