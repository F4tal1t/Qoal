package utils

import (
	"log"

	"github.com/qoal/file-processor/models"
	"gorm.io/gorm"
)

// MigrateDatabase runs all database migrations
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Create users table manually to avoid AutoMigrate issues
	log.Println("Creating users table...")

	// Execute raw SQL to create the table
	sql := `
	CREATE TABLE IF NOT EXISTS qoal_user (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if err := db.Exec(sql).Error; err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
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
