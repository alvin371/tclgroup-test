package dto

import (
	"time"

	"github.com/google/uuid"
)

// ProductOutput is the use-case output for a single product.
type ProductOutput struct {
	ID         uuid.UUID `json:"id"`
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	CustomerID uuid.UUID `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProductListOutput wraps a paginated list of products.
type ProductListOutput struct {
	Items []*ProductOutput `json:"items"`
	Total int64            `json:"total"`
}

// InventoryOutput is the use-case output for a single inventory record.
type InventoryOutput struct {
	ID             uuid.UUID `json:"id"`
	ProductID      uuid.UUID `json:"product_id"`
	ProductSKU     string    `json:"product_sku"`
	ProductName    string    `json:"product_name"`
	PhysicalStock  int64     `json:"physical_stock"`
	Reserved       int64     `json:"reserved"`
	AvailableStock int64     `json:"available_stock"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// InventoryListOutput wraps a paginated list of inventory records.
type InventoryListOutput struct {
	Items []*InventoryOutput `json:"items"`
	Total int64              `json:"total"`
}

// StockInOutput is the use-case output for a single stock-in.
type StockInOutput struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int64     `json:"quantity"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	UnitPrice   float64   `json:"unit_price"`
	PerformedBy string    `json:"performed_by"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockInListOutput wraps a paginated list of stock-ins.
type StockInListOutput struct {
	Items []*StockInOutput `json:"items"`
	Total int64            `json:"total"`
}

// StockOutOutput is the use-case output for a single stock-out.
type StockOutOutput struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int64     `json:"quantity"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	UnitPrice   float64   `json:"unit_price"`
	PerformedBy string    `json:"performed_by"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockOutListOutput wraps a paginated list of stock-outs.
type StockOutListOutput struct {
	Items []*StockOutOutput `json:"items"`
	Total int64             `json:"total"`
}

// ReportOutput is the use-case output for report queries.
type ReportOutput struct {
	Items []*StockInOutput  `json:"stock_in_items,omitempty"`
	Out   []*StockOutOutput `json:"stock_out_items,omitempty"`
	Total int64             `json:"total"`
}

// TimelineEntry is a single entry in a transaction timeline.
type TimelineEntry struct {
	Status      string    `json:"status"`
	Description string    `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// StockInDetailOutput is the detailed use-case output for a single stock-in.
type StockInDetailOutput struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	ProductSKU  string          `json:"product_sku"`
	ProductName string          `json:"product_name"`
	Quantity    int64           `json:"quantity"`
	UnitPrice   float64         `json:"unit_price"`
	Status      string          `json:"status"`
	PerformedBy string          `json:"performed_by"`
	Location    string          `json:"location"`
	Notes       string          `json:"notes"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Timeline    []TimelineEntry `json:"timeline"`
}

// StockOutDetailOutput is the detailed use-case output for a single stock-out.
type StockOutDetailOutput struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	ProductSKU  string          `json:"product_sku"`
	ProductName string          `json:"product_name"`
	Quantity    int64           `json:"quantity"`
	UnitPrice   float64         `json:"unit_price"`
	Status      string          `json:"status"`
	PerformedBy string          `json:"performed_by"`
	Location    string          `json:"location"`
	Notes       string          `json:"notes"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Timeline    []TimelineEntry `json:"timeline"`
}

// HistoryOutput is the use-case output for a single history record.
type HistoryOutput struct {
	ID         uuid.UUID `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Action     string    `json:"action"`
	OldStatus  string    `json:"old_status"`
	NewStatus  string    `json:"new_status"`
	Quantity   int64     `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
}
