package entity

import (
	"time"

	"github.com/google/uuid"
)

// History is an append-only audit log entry for state changes.
type History struct {
	ID         uuid.UUID
	EntityType string
	EntityID   uuid.UUID
	Action     string
	OldStatus  string
	NewStatus  string
	Quantity   int64
	CreatedAt  time.Time
}
