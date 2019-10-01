package models

// Client small wrapper to hold connection parameter
type Client struct {
	Address string
	Port    string
}

// Job required fields for a job
type Job struct {
	JobType string `json:"job_type"`
}
