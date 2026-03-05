package idgen

import "github.com/google/uuid"

// UUIDGenerator implements port.IDGenerator using the google/uuid library.
type UUIDGenerator struct{}

// NewID returns a new random UUID.
func (UUIDGenerator) NewID() uuid.UUID {
	return uuid.New()
}
