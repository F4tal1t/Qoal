package models

type ProcessingJob struct {
	JobID        string
	UserID       int
	InputPath    string
	OutputPath   string
	SourceFormat string
	TargetFormat string
	Type         string
	Status       JobStatus
	Progress     int
	Settings     map[string]interface{}
}
