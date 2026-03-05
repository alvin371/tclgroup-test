package dto

import "github.com/google/uuid"

// CreateProductInput carries data needed to create a product.
type CreateProductInput struct {
	SKU        string
	Name       string
	CustomerID uuid.UUID
}

// CreateStockInInput carries data needed to create a stock-in.
type CreateStockInInput struct {
	ProductID   uuid.UUID
	Quantity    int64
	Notes       string
	UnitPrice   float64
	PerformedBy string
	Location    string
}

// AdvanceStockInInput carries data needed to advance a stock-in status.
type AdvanceStockInInput struct {
	ID        uuid.UUID
	NewStatus string
}

// AllocateStockOutInput carries data needed to allocate stock.
type AllocateStockOutInput struct {
	ProductID   uuid.UUID
	Quantity    int64
	Notes       string
	UnitPrice   float64
	PerformedBy string
	Location    string
}

// ExecuteStockOutInput carries data needed to execute a stock-out transition.
type ExecuteStockOutInput struct {
	ID        uuid.UUID
	NewStatus string
}

// AdjustStockInput carries data needed to adjust inventory stock directly.
type AdjustStockInput struct {
	ProductID uuid.UUID
	NewQty    int64
	Notes     string
}
