package apikeys

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents a stored API credential without exposing the raw secret.
type APIKey struct {
	ID        uuid.UUID  `json:"id"`
	Prefix    string     `json:"prefix"`
	Label     string     `json:"label,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type createRequest struct {
	Label string `json:"label"`
}
