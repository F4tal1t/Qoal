package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qoal/file-processor/config"
	"github.com/qoal/file-processor/handlers"
	"github.com/qoal/file-processor/middleware"
	"github.com/qoal/file-processor/services"
	"github.com/qoal/file-processor/utils"
	"github.com/qoal/file-processor/worker"
)

func main() {
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
	var jobService *services.JobService
	if redisClient != nil {
		jobService = services.NewJobService(redisClient)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	var jobHandler *handlers.JobHandler
	if jobService != nil {
		jobHandler = handlers.NewJobHandler(jobService)
	}
	uploadHandler := handlers.NewUploadHandler()

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth(authService))
	{
		protected.GET("/auth/profile", authHandler.GetProfile)
		if jobHandler != nil {
			protected.POST("/process", jobHandler.CreateJobHandler)
			protected.GET("/status/:id", jobHandler.GetJobStatusHandler)
		}
		protected.POST("/upload", uploadHandler.HandleUpload)
		protected.GET("/uploads", uploadHandler.HandleListUploads)
		protected.DELETE("/uploads/:filename", uploadHandler.HandleDeleteUpload)
	}

	// Start worker in background (only if Redis is available)
	if jobService != nil {
		go func() {
			processor := worker.NewProcessor(jobService, cfg, redisClient)
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
