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
	UserID      uuid.UUID       `json:"user_id,omitempty"`
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

// Decode unmarshals the job payload into the provided destination.
func (j Job) Decode(dest any) error {
	if len(j.Payload) == 0 {
		return errors.New("jobs: empty payload")
	}
	return json.Unmarshal(j.Payload, dest)
}

type enqueueRequest struct {
	Name        string         `json:"name" validate:"required"`
	Payload     map[string]any `json:"payload"`
	RunAt       *time.Time     `json:"run_at"`
	MaxAttempts int            `json:"max_attempts" validate:"omitempty,min=1,max=25"`
}

// ListOptions defines optional filters for job queries.
type ListOptions struct {
	Status string
	Limit  int
	Offset int
	UserID *uuid.UUID
}
