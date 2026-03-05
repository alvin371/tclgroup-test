package inventory

import "errors"

// Domain sentinel errors.
var (
	ErrInvalidQuantity         = errors.New("quantity must be positive")
	ErrInsufficientStock       = errors.New("insufficient available stock")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrCannotCancelDoneStockIn = errors.New("cannot cancel completed stock in")
	ErrProductNotFound         = errors.New("product not found")
	ErrInventoryNotFound       = errors.New("inventory not found")
	ErrStockInNotFound         = errors.New("stock in not found")
	ErrStockOutNotFound        = errors.New("stock out not found")
	ErrReservationNotFound     = errors.New("reservation not found")
)
