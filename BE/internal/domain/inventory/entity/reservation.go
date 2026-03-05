package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
)

// Reservation holds a quantity of stock reserved for a pending StockOut.
type Reservation struct {
	ID         uuid.UUID
	StockOutID uuid.UUID
	ProductID  uuid.UUID
	Quantity   int64
	Status     valueobject.ReservationStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Consume transitions the reservation from ACTIVE to CONSUMED.
func (r *Reservation) Consume() error {
	if r.Status != valueobject.ReservationActive {
		return inventory.ErrInvalidStatusTransition
	}
	r.Status = valueobject.ReservationConsumed
	return nil
}

// Release transitions the reservation from ACTIVE to RELEASED.
func (r *Reservation) Release() error {
	if r.Status != valueobject.ReservationActive {
		return inventory.ErrInvalidStatusTransition
	}
	r.Status = valueobject.ReservationReleased
	return nil
}
