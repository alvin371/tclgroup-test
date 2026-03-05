package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/tclgroup/stock-management/internal/domain/inventory"
)

// Inventory holds the physical stock for a product.
type Inventory struct {
	ID            uuid.UUID
	ProductID     uuid.UUID
	PhysicalStock int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// AvailableStock computes the available (non-reserved) stock.
func (inv *Inventory) AvailableStock(reserved int64) int64 {
	return inv.PhysicalStock - reserved
}

// AddStock increases the physical stock when a StockIn is DONE.
func (inv *Inventory) AddStock(qty int64) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	inv.PhysicalStock += qty
	return nil
}

// DeductStock decreases the physical stock when a StockOut is DONE.
func (inv *Inventory) DeductStock(qty int64) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	if inv.PhysicalStock < qty {
		return inventory.ErrInsufficientStock
	}
	inv.PhysicalStock -= qty
	return nil
}

// Adjust sets the physical stock to a new absolute value.
func (inv *Inventory) Adjust(newQty int64) error {
	if newQty < 0 {
		return inventory.ErrInvalidQuantity
	}
	inv.PhysicalStock = newQty
	return nil
}
