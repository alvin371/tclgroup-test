package dto

import "github.com/google/uuid"

// CreateProductRequest is the HTTP request body for product creation.
type CreateProductRequest struct {
	SKU        string    `json:"sku"         binding:"required"`
	Name       string    `json:"name"        binding:"required"`
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`
}

// CreateStockInRequest is the HTTP request body for stock-in creation.
type CreateStockInRequest struct {
	ProductID   uuid.UUID `json:"product_id"   binding:"required"`
	Quantity    int64     `json:"quantity"     binding:"required,gt=0"`
	Notes       string    `json:"notes"`
	UnitPrice   float64   `json:"unit_price"   binding:"min=0"`
	PerformedBy string    `json:"performed_by"`
	Location    string    `json:"location"`
}

// AdvanceStockInRequest is the HTTP request body for advancing a stock-in.
type AdvanceStockInRequest struct {
	Status string `json:"status" binding:"required"`
}

// AllocateStockOutRequest is the HTTP request body for allocating stock.
type AllocateStockOutRequest struct {
	ProductID   uuid.UUID `json:"product_id"   binding:"required"`
	Quantity    int64     `json:"quantity"     binding:"required,gt=0"`
	Notes       string    `json:"notes"`
	UnitPrice   float64   `json:"unit_price"   binding:"min=0"`
	PerformedBy string    `json:"performed_by"`
	Location    string    `json:"location"`
}

// ExecuteStockOutRequest is the HTTP request body for executing a stock-out transition.
type ExecuteStockOutRequest struct {
	Status string `json:"status" binding:"required"`
}

// AdjustStockRequest is the HTTP request body for adjusting inventory.
type AdjustStockRequest struct {
	NewQty int64  `json:"new_qty" binding:"min=0"`
	Notes  string `json:"notes"`
}

// PaginationQuery carries query parameters for pagination.
type PaginationQuery struct {
	Page    int `form:"page"     binding:"omitempty,min=1"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}
