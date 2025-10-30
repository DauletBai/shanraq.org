package jobs

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNoJobs indicates no runnable jobs are currently available.
var ErrNoJobs = errors.New("no jobs available")

// Job mirrors the job_queue schema.
type Job struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Payload     json.RawMessage `json:"payload"`
	RunAt       time.Time       `json:"run_at"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	Status      string          `json:"status"`
	LastError   *string         `json:"last_error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type enqueueRequest struct {
	Name        string         `json:"name"`
	Payload     map[string]any `json:"payload"`
	RunAt       *time.Time     `json:"run_at"`
	MaxAttempts int            `json:"max_attempts"`
}

// ListOptions defines optional filters for job queries.
type ListOptions struct {
	Status string
	Limit  int
	Offset int
}
