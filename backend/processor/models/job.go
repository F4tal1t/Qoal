package models

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

// Job represents a file conversion job in the database
type Job struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	JobID            string     `gorm:"unique;not null" json:"job_id"`     // UUID string for external reference
	UserID           string     `gorm:"not null" json:"user_id"`           // Foreign key to User (UUID string)
	OriginalFilename string     `gorm:"not null" json:"original_filename"` // Original file name
	FileSize         int64      `gorm:"not null" json:"file_size"`         // File size in bytes
	SourceFormat     string     `gorm:"not null" json:"source_format"`     // Source file extension
	TargetFormat     string     `gorm:"not null" json:"target_format"`     // Target file extension
	Status           string     `gorm:"default:'pending'" json:"status"`   // Job status
	InputPath        string     `gorm:"not null" json:"input_path"`        // Local input file path
	OutputPath       string     `json:"output_path"`                       // Local output file path (empty until completed)
	Error            string     `json:"error,omitempty"`                   // Error message if failed
	CompletedAt      *time.Time `json:"completed_at,omitempty"`            // Completion timestamp (null until completed)
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TableName specifies the custom table name for Job model
func (Job) TableName() string {
	return "qoal_job"
}
