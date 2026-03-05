package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/tclgroup/stock-management/internal/domain/inventory/entity"
)

// Pagination carries page and limit parameters.
type Pagination struct {
	Page  int
	Limit int
}

// ProductFilter carries filter criteria for product queries.
type ProductFilter struct {
	CustomerID *uuid.UUID
	SKU        *string
	Name       *string
}

// InventoryFilter carries filter criteria for inventory queries.
type InventoryFilter struct {
	ProductID  *uuid.UUID
	CustomerID *uuid.UUID
}

// InventoryRow is a read model joining inventory with product data.
type InventoryRow struct {
	Inventory *entity.Inventory
	Product   *entity.Product
	Reserved  int64
}

// StockInFilter carries filter criteria for stock-in queries.
type StockInFilter struct {
	ProductID *uuid.UUID
	Status    *string
}

// StockOutFilter carries filter criteria for stock-out queries.
type StockOutFilter struct {
	ProductID *uuid.UUID
	Status    *string
}

// ProductRepo is the repository interface for products.
type ProductRepo interface {
	Create(ctx context.Context, p *entity.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	FindAll(ctx context.Context, filter ProductFilter, page Pagination) ([]*entity.Product, int64, error)
}

// InventoryRepo is the repository interface for inventories.
type InventoryRepo interface {
	FindByProductID(ctx context.Context, productID uuid.UUID) (*entity.Inventory, error)
	FindAllWithFilter(ctx context.Context, filter InventoryFilter, page Pagination) ([]*InventoryRow, int64, error)
	Update(ctx context.Context, inv *entity.Inventory) error
	GetTotalReserved(ctx context.Context, productID uuid.UUID) (int64, error)
	LockByProductID(ctx context.Context, productID uuid.UUID) (*entity.Inventory, error)
	Create(ctx context.Context, inv *entity.Inventory) error
}

// StockInRepo is the repository interface for stock-in transactions.
type StockInRepo interface {
	Create(ctx context.Context, s *entity.StockIn) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.StockIn, error)
	Update(ctx context.Context, s *entity.StockIn) error
	FindAll(ctx context.Context, filter StockInFilter, page Pagination) ([]*entity.StockIn, int64, error)
}

// StockOutRepo is the repository interface for stock-out transactions.
type StockOutRepo interface {
	Create(ctx context.Context, s *entity.StockOut) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.StockOut, error)
	Update(ctx context.Context, s *entity.StockOut) error
	FindAll(ctx context.Context, filter StockOutFilter, page Pagination) ([]*entity.StockOut, int64, error)
}

// ReservationRepo is the repository interface for reservations.
type ReservationRepo interface {
	Create(ctx context.Context, r *entity.Reservation) error
	FindByStockOutID(ctx context.Context, stockOutID uuid.UUID) (*entity.Reservation, error)
	Update(ctx context.Context, r *entity.Reservation) error
}

// HistoryRepo is the repository interface for audit histories.
type HistoryRepo interface {
	Create(ctx context.Context, h *entity.History) error
	FindByEntityID(ctx context.Context, entityID uuid.UUID) ([]*entity.History, error)
}
