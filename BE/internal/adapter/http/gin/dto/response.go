package dto

import (
	"time"

	"github.com/google/uuid"
)

// ProductResponse is the HTTP response for a product.
type ProductResponse struct {
	ID         uuid.UUID `json:"id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// InventoryResponse is the HTTP response for an inventory record.
type InventoryResponse struct {
	ID             uuid.UUID `json:"id"`
	ProductID      uuid.UUID `json:"product_id"`
	ProductSKU     string    `json:"product_sku"`
	ProductName    string    `json:"product_name"`
	PhysicalStock  int64     `json:"physical_stock"`
	Reserved       int64     `json:"reserved"`
	AvailableStock int64     `json:"available_stock"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// StockInResponse is the HTTP response for a stock-in record.
type StockInResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StockOutResponse is the HTTP response for a stock-out record.
type StockOutResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Status    string    `json:"status"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
