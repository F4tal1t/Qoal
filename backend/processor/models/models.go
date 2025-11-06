package models

type ProcessingJob struct {
	JobID        string                 `json:"job_id"`
	UserID       string                 `json:"user_id"`
	InputPath    string                 `json:"input_path"`
	OutputPath   string                 `json:"output_path"`
	SourceFormat string                 `json:"source_format"`
	TargetFormat string                 `json:"target_format"`
	Type         string                 `json:"type"`
	Status       JobStatus              `json:"status"`
	Progress     int                    `json:"progress"`
	Settings     map[string]interface{} `json:"settings"`
}
