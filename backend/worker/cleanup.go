package worker

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"

	"github.com/qoal/file-processor/models"
)

type CleanupWorker struct {
	db           *gorm.DB
	uploadDir    string
	processedDir string
	tempDir      string
}

func NewCleanupWorker(db *gorm.DB, uploadDir, processedDir, tempDir string) *CleanupWorker {
	return &CleanupWorker{
		db:           db,
		uploadDir:    uploadDir,
		processedDir: processedDir,
		tempDir:      tempDir,
	}
}

func (w *CleanupWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Println("Cleanup worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup worker stopped")
			return
		case <-ticker.C:
			w.cleanupOldFiles()
		}
	}
}

func (w *CleanupWorker) cleanupOldFiles() {
	cutoffTime := time.Now().Add(-24 * time.Hour)

	// Clean completed jobs older than 24 hours
	var oldJobs []models.Job
	w.db.Where("status = ? AND updated_at < ?", models.StatusCompleted, cutoffTime).Find(&oldJobs)

	for _, job := range oldJobs {
		if job.InputPath != "" {
			os.Remove(job.InputPath)
		}
		if job.OutputPath != "" {
			os.Remove(job.OutputPath)
		}
		w.db.Delete(&job)
	}

	// Clean temp files older than 1 hour
	tempCutoff := time.Now().Add(-1 * time.Hour)
	filepath.Walk(w.tempDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.ModTime().Before(tempCutoff) {
			os.Remove(path)
		}
		return nil
	})

	log.Printf("Cleanup completed: removed %d old jobs", len(oldJobs))
}
