package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
)

// StockOut represents an outbound stock transaction.
type StockOut struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	Quantity    int64
	Status      valueobject.StockOutStatus
	Notes       string
	UnitPrice   float64
	PerformedBy string
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Advance transitions the StockOut to the given next status.
func (s *StockOut) Advance(next valueobject.StockOutStatus) error {
	if !s.Status.CanTransitionTo(next) {
		return inventory.ErrInvalidStatusTransition
	}
	s.Status = next
	return nil
}

// Cancel transitions the StockOut to CANCELLED.
// Returns an error if the transition is not allowed.
func (s *StockOut) Cancel() error {
	return s.Advance(valueobject.StockOutCancelled)
}

// NeedsRollback returns true if cancelling this StockOut requires releasing a reservation.
// This is the case when the StockOut was already IN_PROGRESS.
func (s *StockOut) NeedsRollback() bool {
	return s.Status == valueobject.StockOutInProgress
}
