package inventory

import (
	"context"
	"fmt"

	domain "github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/entity"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// InventoryUseCase handles inventory read and adjustment operations.
type InventoryUseCase struct {
	inventoryRepo port.InventoryRepo
	productRepo   port.ProductRepo
	historyRepo   port.HistoryRepo
	txManager     port.TxManager
	clock         port.Clock
	idgen         port.IDGenerator
}

// NewInventoryUseCase creates a new InventoryUseCase.
func NewInventoryUseCase(
	inventoryRepo port.InventoryRepo,
	productRepo port.ProductRepo,
	historyRepo port.HistoryRepo,
	txManager port.TxManager,
	clock port.Clock,
	idgen port.IDGenerator,
) *InventoryUseCase {
	return &InventoryUseCase{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		historyRepo:   historyRepo,
		txManager:     txManager,
		clock:         clock,
		idgen:         idgen,
	}
}

// List returns a paginated list of inventory records with computed available stock.
func (uc *InventoryUseCase) List(ctx context.Context, filter port.InventoryFilter, page port.Pagination) (*dto.InventoryListOutput, error) {
	rows, total, err := uc.inventoryRepo.FindAllWithFilter(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("list inventory: %w", err)
	}

	items := make([]*dto.InventoryOutput, len(rows))
	for i, row := range rows {
		items[i] = &dto.InventoryOutput{
			ID:             row.Inventory.ID,
			ProductID:      row.Inventory.ProductID,
			ProductSKU:     row.Product.SKU,
			ProductName:    row.Product.Name,
			PhysicalStock:  row.Inventory.PhysicalStock,
			Reserved:       row.Reserved,
			AvailableStock: row.Inventory.AvailableStock(row.Reserved),
			UpdatedAt:      row.Inventory.UpdatedAt,
		}
	}

	return &dto.InventoryListOutput{Items: items, Total: total}, nil
}

// Adjust directly sets the physical stock to a new value and logs the change.
func (uc *InventoryUseCase) Adjust(ctx context.Context, input dto.AdjustStockInput) error {
	return uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		inv, err := uc.inventoryRepo.LockByProductID(txCtx, input.ProductID)
		if err != nil {
			return err
		}

		oldQty := inv.PhysicalStock
		if err = inv.Adjust(input.NewQty); err != nil {
			return err
		}
		inv.UpdatedAt = uc.clock.Now()

		if err = uc.inventoryRepo.Update(txCtx, inv); err != nil {
			return fmt.Errorf("update inventory: %w", err)
		}

		h := &entity.History{
			ID:         uc.idgen.NewID(),
			EntityType: "inventory",
			EntityID:   inv.ID,
			Action:     "adjust",
			OldStatus:  fmt.Sprintf("%d", oldQty),
			NewStatus:  fmt.Sprintf("%d", input.NewQty),
			Quantity:   input.NewQty,
			CreatedAt:  uc.clock.Now(),
		}
		return uc.historyRepo.Create(txCtx, h)
	})
}

// CreateProduct creates a product and initialises its inventory record.
func (uc *InventoryUseCase) CreateProduct(ctx context.Context, input dto.CreateProductInput) (*dto.ProductOutput, error) {
	if input.SKU == "" {
		return nil, domain.ErrInvalidQuantity // reuse or add SKU error; kept minimal
	}

	var result *dto.ProductOutput

	err := uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		now := uc.clock.Now()
		p := &entity.Product{
			ID:         uc.idgen.NewID(),
			SKU:        input.SKU,
			Name:       input.Name,
			CustomerID: input.CustomerID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := uc.productRepo.Create(txCtx, p); err != nil {
			return fmt.Errorf("create product: %w", err)
		}

		inv := &entity.Inventory{
			ID:            uc.idgen.NewID(),
			ProductID:     p.ID,
			PhysicalStock: 0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := uc.inventoryRepo.Create(txCtx, inv); err != nil {
			return fmt.Errorf("create inventory: %w", err)
		}

		result = &dto.ProductOutput{
			ID:         p.ID,
			SKU:        p.SKU,
			Name:       p.Name,
			CustomerID: p.CustomerID,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListProducts returns a paginated list of products.
func (uc *InventoryUseCase) ListProducts(ctx context.Context, filter port.ProductFilter, page port.Pagination) (*dto.ProductListOutput, error) {
	items, total, err := uc.productRepo.FindAll(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	out := make([]*dto.ProductOutput, len(items))
	for i, p := range items {
		out[i] = &dto.ProductOutput{
			ID:         p.ID,
			SKU:        p.SKU,
			Name:       p.Name,
			CustomerID: p.CustomerID,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		}
	}
	return &dto.ProductListOutput{Items: out, Total: total}, nil
}
