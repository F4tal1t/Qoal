package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/utils"
)

// 🎓 **Go Tutorial: Simple Testing Without Framework**
// This is a simple way to test our code without external testing frameworks
func main() {
	fmt.Println("Testing Qoal Backend Components")
	fmt.Println("=====================================")

	// Test 1: Document Extensions
	fmt.Println("\n Testing Document Extensions:")
	testDocumentExtensions()

	// Test 2: Video Extensions
	fmt.Println("\n Testing Video Extensions:")
	testVideoExtensions()

	// Test 3: Job Model
	fmt.Println("\n Testing Job Model:")
	testJobModel()

	// Test 4: File Operations
	fmt.Println("\n Testing File Operations:")
	testFileOperations()

	fmt.Println("\n All tests completed!")
}

func testDocumentExtensions() {
	formats := []string{"pdf", "docx", "doc", "txt", "rtf", "odt", "invalid"}

	for _, format := range formats {
		ext, err := utils.GetDocumentExtension(format)
		if err != nil {
			fmt.Printf("   %s: %s\n", format, err.Error())
		} else {
			fmt.Printf("   %s -> %s\n", format, ext)
		}
	}
}

func testVideoExtensions() {
	formats := []string{"mp4", "avi", "mov", "mkv", "webm", "invalid"}

	for _, format := range formats {
		ext, err := utils.GetVideoExtension(format)
		if err != nil {
			fmt.Printf("   %s: %s\n", format, err.Error())
		} else {
			fmt.Printf("   %s -> %s\n", format, ext)
		}
	}
}

func testJobModel() {
	job := models.Job{
		ID:        "test-job-123",
		UserID:    "user-456",
		InputFile: "/uploads/test.pdf",
		Status:    models.StatusPending,
		CreatedAt: 1640995200, // Unix timestamp
	}

	fmt.Printf("   Job ID: %s\n", job.ID)
	fmt.Printf("   User ID: %s\n", job.UserID)
	fmt.Printf("   Input File: %s\n", job.InputFile)
	fmt.Printf("   Status: %s\n", job.Status)
	fmt.Printf("   Created: %d\n", job.CreatedAt)

	// Test status changes
	job.Status = models.StatusProcessing
	fmt.Printf("   Status changed to: %s\n", job.Status)

	job.Status = models.StatusCompleted
	job.OutputFile = "/output/test_converted.docx"
	job.CompletedAt = 1640995260
	fmt.Printf("   Status changed to: %s\n", job.Status)
	fmt.Printf("   Output File: %s\n", job.OutputFile)
	fmt.Printf("   Completed: %d\n", job.CompletedAt)
}

func testFileOperations() {
	// Test creating uploads directory
	uploadDir := "test_uploads"
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Printf("   Failed to create directory: %v", err)
		return
	}
	fmt.Printf("   Created directory: %s\n", uploadDir)

	// Test creating a test file
	testFile := filepath.Join(uploadDir, "test.txt")
	content := []byte("Hello, Qoal! This is a test file.")
	err = os.WriteFile(testFile, content, 0644)
	if err != nil {
		log.Printf("   Failed to write file: %v", err)
		return
	}
	fmt.Printf("   Created test file: %s\n", testFile)

	// Test reading file info
	info, err := os.Stat(testFile)
	if err != nil {
		log.Printf("   Failed to stat file: %v", err)
		return
	}
	fmt.Printf("   File size: %d bytes\n", info.Size())
	fmt.Printf("   Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))

	// Test reading file content
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		log.Printf("   Failed to read file: %v", err)
		return
	}
	fmt.Printf("   File content: %s\n", string(readContent))

	// Cleanup
	err = os.Remove(testFile)
	if err != nil {
		log.Printf("   Failed to remove file: %v", err)
		return
	}
	fmt.Printf("    Removed test file\n")

	err = os.Remove(uploadDir)
	if err != nil {
		log.Printf("   Failed to remove directory: %v", err)
		return
	}
	fmt.Printf("    Removed test directory\n")
}