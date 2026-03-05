package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domain "github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/entity"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// StockInUseCase handles stock-in business operations.
type StockInUseCase struct {
	stockInRepo   port.StockInRepo
	inventoryRepo port.InventoryRepo
	historyRepo   port.HistoryRepo
	productRepo   port.ProductRepo
	txManager     port.TxManager
	clock         port.Clock
	idgen         port.IDGenerator
}

// NewStockInUseCase creates a new StockInUseCase.
func NewStockInUseCase(
	stockInRepo port.StockInRepo,
	inventoryRepo port.InventoryRepo,
	historyRepo port.HistoryRepo,
	txManager port.TxManager,
	clock port.Clock,
	idgen port.IDGenerator,
	productRepo port.ProductRepo,
) *StockInUseCase {
	return &StockInUseCase{
		stockInRepo:   stockInRepo,
		inventoryRepo: inventoryRepo,
		historyRepo:   historyRepo,
		productRepo:   productRepo,
		txManager:     txManager,
		clock:         clock,
		idgen:         idgen,
	}
}

// Create creates a new stock-in record in CREATED status.
func (uc *StockInUseCase) Create(ctx context.Context, input dto.CreateStockInInput) (*dto.StockInOutput, error) {
	if input.Quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	now := uc.clock.Now()
	s := &entity.StockIn{
		ID:          uc.idgen.NewID(),
		ProductID:   input.ProductID,
		Quantity:    input.Quantity,
		Status:      valueobject.StockInCreated,
		Notes:       input.Notes,
		UnitPrice:   input.UnitPrice,
		PerformedBy: input.PerformedBy,
		Location:    input.Location,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.stockInRepo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("create stock in: %w", err)
	}

	return mapStockInOutput(s), nil
}

// Advance transitions a stock-in to the given status.
// When status is DONE, it locks the inventory and adds the stock.
func (uc *StockInUseCase) Advance(ctx context.Context, id uuid.UUID, newStatus valueobject.StockInStatus) (*dto.StockInOutput, error) {
	var result *dto.StockInOutput

	err := uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		s, err := uc.stockInRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		oldStatus := s.Status
		if err = s.Advance(newStatus); err != nil {
			return err
		}
		s.UpdatedAt = uc.clock.Now()

		if newStatus == valueobject.StockInDone {
			inv, err := uc.inventoryRepo.LockByProductID(txCtx, s.ProductID)
			if err != nil {
				return err
			}
			if err = inv.AddStock(s.Quantity); err != nil {
				return err
			}
			inv.UpdatedAt = uc.clock.Now()
			if err = uc.inventoryRepo.Update(txCtx, inv); err != nil {
				return fmt.Errorf("update inventory: %w", err)
			}
		}

		if err = uc.stockInRepo.Update(txCtx, s); err != nil {
			return fmt.Errorf("update stock in: %w", err)
		}

		h := &entity.History{
			ID:         uc.idgen.NewID(),
			EntityType: "stock_in",
			EntityID:   s.ID,
			Action:     "advance",
			OldStatus:  string(oldStatus),
			NewStatus:  string(newStatus),
			Quantity:   s.Quantity,
			CreatedAt:  uc.clock.Now(),
		}
		if err = uc.historyRepo.Create(txCtx, h); err != nil {
			return fmt.Errorf("create history: %w", err)
		}

		result = mapStockInOutput(s)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Cancel cancels a stock-in (must not be in DONE status).
func (uc *StockInUseCase) Cancel(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		s, err := uc.stockInRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		oldStatus := s.Status
		if err = s.Cancel(); err != nil {
			return err
		}
		s.UpdatedAt = uc.clock.Now()

		if err = uc.stockInRepo.Update(txCtx, s); err != nil {
			return fmt.Errorf("update stock in: %w", err)
		}

		h := &entity.History{
			ID:         uc.idgen.NewID(),
			EntityType: "stock_in",
			EntityID:   s.ID,
			Action:     "cancel",
			OldStatus:  string(oldStatus),
			NewStatus:  string(valueobject.StockInCancelled),
			Quantity:   s.Quantity,
			CreatedAt:  uc.clock.Now(),
		}
		return uc.historyRepo.Create(txCtx, h)
	})
}

// List returns a paginated list of stock-ins.
func (uc *StockInUseCase) List(ctx context.Context, filter port.StockInFilter, page port.Pagination) (*dto.StockInListOutput, error) {
	items, total, err := uc.stockInRepo.FindAll(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("list stock ins: %w", err)
	}

	out := make([]*dto.StockInOutput, len(items))
	for i, s := range items {
		out[i] = mapStockInOutput(s)
	}
	return &dto.StockInListOutput{Items: out, Total: total}, nil
}

func mapStockInOutput(s *entity.StockIn) *dto.StockInOutput {
	return &dto.StockInOutput{
		ID:          s.ID,
		ProductID:   s.ProductID,
		Quantity:    s.Quantity,
		Status:      string(s.Status),
		Notes:       s.Notes,
		UnitPrice:   s.UnitPrice,
		PerformedBy: s.PerformedBy,
		Location:    s.Location,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// GetDetail retrieves a detailed stock-in record including product info and timeline.
func (uc *StockInUseCase) GetDetail(ctx context.Context, id uuid.UUID) (*dto.StockInDetailOutput, error) {
	s, err := uc.stockInRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product, err := uc.productRepo.FindByID(ctx, s.ProductID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	histories, err := uc.historyRepo.FindByEntityID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get histories: %w", err)
	}

	return &dto.StockInDetailOutput{
		ID:          s.ID,
		ProductID:   s.ProductID,
		ProductSKU:  product.SKU,
		ProductName: product.Name,
		Quantity:    s.Quantity,
		UnitPrice:   s.UnitPrice,
		Status:      string(s.Status),
		PerformedBy: s.PerformedBy,
		Location:    s.Location,
		Notes:       s.Notes,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Timeline:    buildStockInTimeline(s, histories),
	}, nil
}

func buildStockInTimeline(s *entity.StockIn, histories []*entity.History) []dto.TimelineEntry {
	entries := []dto.TimelineEntry{{
		Status:      "CREATED",
		Description: "Transaction record created by System API.",
		OccurredAt:  s.CreatedAt,
	}}
	for _, h := range histories {
		entries = append(entries, dto.TimelineEntry{
			Status:      h.NewStatus,
			Description: stockInStatusDescription(h.NewStatus, h.OldStatus, s.PerformedBy),
			OccurredAt:  h.CreatedAt,
		})
	}
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries
}

func stockInStatusDescription(newStatus, oldStatus, actor string) string {
	if actor == "" {
		actor = "System"
	}
	switch newStatus {
	case "IN_PROGRESS":
		return "Items being verified and received at the loading dock."
	case "DONE":
		return fmt.Sprintf("Transaction completed and stock updated in the system by %s.", actor)
	case "CANCELLED":
		return fmt.Sprintf("Transaction cancelled from %s by %s.", oldStatus, actor)
	default:
		return fmt.Sprintf("Status changed from %s to %s.", oldStatus, newStatus)
	}
}
