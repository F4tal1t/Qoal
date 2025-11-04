package models

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	InputFile   string    `json:"input_file"`
	OutputFile  string    `json:"output_file,omitempty"`
	Status      JobStatus `json:"status"`
	CreatedAt   int64     `json:"created_at"`
	CompletedAt int64     `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}
