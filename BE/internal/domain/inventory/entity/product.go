package entity

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the inventory system.
type Product struct {
	ID         uuid.UUID
	SKU        string
	Name       string
	CustomerID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
