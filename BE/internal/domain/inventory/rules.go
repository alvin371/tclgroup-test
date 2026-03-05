package inventory

import "github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"

// ValidateStockInTransition returns an error if the transition is not allowed.
func ValidateStockInTransition(from, to valueobject.StockInStatus) error {
	if !from.CanTransitionTo(to) {
		return ErrInvalidStatusTransition
	}
	return nil
}

// ValidateStockOutTransition returns an error if the transition is not allowed.
func ValidateStockOutTransition(from, to valueobject.StockOutStatus) error {
	if !from.CanTransitionTo(to) {
		return ErrInvalidStatusTransition
	}
	return nil
}

// CanAllocate returns true if the requested quantity can be satisfied given
// the current physical stock and the total reserved quantity.
func CanAllocate(physical, reserved, requested int64) bool {
	return physical-reserved >= requested
}
