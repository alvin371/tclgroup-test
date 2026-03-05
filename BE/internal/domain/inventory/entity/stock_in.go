package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
)

// StockIn represents an inbound stock transaction.
type StockIn struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	Quantity    int64
	Status      valueobject.StockInStatus
	Notes       string
	UnitPrice   float64
	PerformedBy string
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Advance transitions the StockIn to the given next status.
func (s *StockIn) Advance(next valueobject.StockInStatus) error {
	if !s.Status.CanTransitionTo(next) {
		return inventory.ErrInvalidStatusTransition
	}
	s.Status = next
	return nil
}

// Cancel attempts to cancel the StockIn. Returns ErrCannotCancelDoneStockIn if already DONE.
func (s *StockIn) Cancel() error {
	if s.Status == valueobject.StockInDone {
		return inventory.ErrCannotCancelDoneStockIn
	}
	return s.Advance(valueobject.StockInCancelled)
}
