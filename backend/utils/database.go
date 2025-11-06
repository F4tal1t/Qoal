package utils

import (
	"fmt"
	"log"

	"github.com/qoal/file-processor/models"
	"gorm.io/gorm"
)

// MigrateDatabase runs all database migrations
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Try to create tables manually if AutoMigrate fails
	// Create User table
	err := db.Exec(`CREATE TABLE IF NOT EXISTS qoal_user (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`).Error
	if err != nil {
		return fmt.Errorf("failed to create User table: %w", err)
	}

	// Drop existing Job table if it exists with old schema
	log.Println("Dropping existing Job table if it exists...")
	db.Exec(`DROP TABLE IF EXISTS qoal_job`)

	// Create Job table with correct schema
	err = db.Exec(`CREATE TABLE IF NOT EXISTS qoal_job (
		id SERIAL PRIMARY KEY,
		job_id VARCHAR(255) UNIQUE NOT NULL,
		user_id UUID NOT NULL,
		original_filename VARCHAR(255) NOT NULL,
		file_size BIGINT NOT NULL,
		source_format VARCHAR(50) NOT NULL,
		target_format VARCHAR(50) NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		input_path TEXT NOT NULL,
		output_path TEXT,
		error TEXT,
		completed_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`).Error
	if err != nil {
		return fmt.Errorf("failed to create Job table: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedDatabase adds initial data if needed
func SeedDatabase(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Check if admin user exists
	var adminUser models.User
	err := db.Table("qoal_user").Where("email = ?", "admin@qoal.com").First(&adminUser).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		// Create admin user
		adminUser = models.User{
			Email:    "admin@qoal.com",
			Password: "admin123", // This will be hashed by the auth service
			Name:     "Admin User",
		}

		// Create user - omit ID to let PostgreSQL generate UUID automatically
		if err := db.Table("qoal_user").Select("email", "password", "name", "created_at", "updated_at").Create(&adminUser).Error; err != nil {
			return err
		}

		log.Println("Admin user created successfully")
	}

	log.Println("Database seeding completed")
	return nil
}
