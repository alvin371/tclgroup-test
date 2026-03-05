package port

import "github.com/google/uuid"

// IDGenerator abstracts UUID generation for testability.
type IDGenerator interface {
	NewID() uuid.UUID
}
