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

// StockOutUseCase handles stock-out business operations with Two-Phase Commitment.
type StockOutUseCase struct {
	stockOutRepo    port.StockOutRepo
	inventoryRepo   port.InventoryRepo
	reservationRepo port.ReservationRepo
	historyRepo     port.HistoryRepo
	productRepo     port.ProductRepo
	txManager       port.TxManager
	clock           port.Clock
	idgen           port.IDGenerator
}

// NewStockOutUseCase creates a new StockOutUseCase.
func NewStockOutUseCase(
	stockOutRepo port.StockOutRepo,
	inventoryRepo port.InventoryRepo,
	reservationRepo port.ReservationRepo,
	historyRepo port.HistoryRepo,
	txManager port.TxManager,
	clock port.Clock,
	idgen port.IDGenerator,
	productRepo port.ProductRepo,
) *StockOutUseCase {
	return &StockOutUseCase{
		stockOutRepo:    stockOutRepo,
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
		historyRepo:     historyRepo,
		productRepo:     productRepo,
		txManager:       txManager,
		clock:           clock,
		idgen:           idgen,
	}
}

// Allocate implements Phase 1 of Two-Phase Commitment:
// locks inventory, checks availability, creates DRAFT StockOut + ACTIVE Reservation.
func (uc *StockOutUseCase) Allocate(ctx context.Context, input dto.AllocateStockOutInput) (*dto.StockOutOutput, error) {
	if input.Quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	var result *dto.StockOutOutput

	err := uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		// Phase 1: lock inventory row
		inv, err := uc.inventoryRepo.LockByProductID(txCtx, input.ProductID)
		if err != nil {
			return err
		}

		// Compute total active reservations
		reserved, err := uc.inventoryRepo.GetTotalReserved(txCtx, input.ProductID)
		if err != nil {
			return fmt.Errorf("get total reserved: %w", err)
		}

		if !domain.CanAllocate(inv.PhysicalStock, reserved, input.Quantity) {
			return domain.ErrInsufficientStock
		}

		now := uc.clock.Now()
		s := &entity.StockOut{
			ID:          uc.idgen.NewID(),
			ProductID:   input.ProductID,
			Quantity:    input.Quantity,
			Status:      valueobject.StockOutDraft,
			Notes:       input.Notes,
			UnitPrice:   input.UnitPrice,
			PerformedBy: input.PerformedBy,
			Location:    input.Location,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err = uc.stockOutRepo.Create(txCtx, s); err != nil {
			return fmt.Errorf("create stock out: %w", err)
		}

		r := &entity.Reservation{
			ID:         uc.idgen.NewID(),
			StockOutID: s.ID,
			ProductID:  input.ProductID,
			Quantity:   input.Quantity,
			Status:     valueobject.ReservationActive,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err = uc.reservationRepo.Create(txCtx, r); err != nil {
			return fmt.Errorf("create reservation: %w", err)
		}

		result = mapStockOutOutput(s)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Execute advances a stock-out to the given status (Phase 2).
func (uc *StockOutUseCase) Execute(ctx context.Context, id uuid.UUID, newStatus valueobject.StockOutStatus) (*dto.StockOutOutput, error) {
	var result *dto.StockOutOutput

	err := uc.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		s, err := uc.stockOutRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		oldStatus := s.Status

		switch newStatus {
		case valueobject.StockOutDone:
			if err = s.Advance(valueobject.StockOutDone); err != nil {
				return err
			}
			s.UpdatedAt = uc.clock.Now()

			// Consume reservation
			res, err := uc.reservationRepo.FindByStockOutID(txCtx, s.ID)
			if err != nil {
				return err
			}
			if err = res.Consume(); err != nil {
				return err
			}
			res.UpdatedAt = uc.clock.Now()
			if err = uc.reservationRepo.Update(txCtx, res); err != nil {
				return fmt.Errorf("update reservation: %w", err)
			}

			// Deduct physical stock (lock first)
			inv, err := uc.inventoryRepo.LockByProductID(txCtx, s.ProductID)
			if err != nil {
				return err
			}
			if err = inv.DeductStock(s.Quantity); err != nil {
				return err
			}
			inv.UpdatedAt = uc.clock.Now()
			if err = uc.inventoryRepo.Update(txCtx, inv); err != nil {
				return fmt.Errorf("update inventory: %w", err)
			}

		case valueobject.StockOutCancelled:
			// Release reservation if cancelling from any valid state
			res, err := uc.reservationRepo.FindByStockOutID(txCtx, s.ID)
			if err != nil {
				return err
			}
			if err = res.Release(); err != nil {
				return err
			}
			res.UpdatedAt = uc.clock.Now()
			if err = uc.reservationRepo.Update(txCtx, res); err != nil {
				return fmt.Errorf("update reservation: %w", err)
			}

			if err = s.Cancel(); err != nil {
				return err
			}
			s.UpdatedAt = uc.clock.Now()

		default:
			if err = s.Advance(newStatus); err != nil {
				return err
			}
			s.UpdatedAt = uc.clock.Now()
		}

		if err = uc.stockOutRepo.Update(txCtx, s); err != nil {
			return fmt.Errorf("update stock out: %w", err)
		}

		h := &entity.History{
			ID:         uc.idgen.NewID(),
			EntityType: "stock_out",
			EntityID:   s.ID,
			Action:     "execute",
			OldStatus:  string(oldStatus),
			NewStatus:  string(s.Status),
			Quantity:   s.Quantity,
			CreatedAt:  uc.clock.Now(),
		}
		if err = uc.historyRepo.Create(txCtx, h); err != nil {
			return fmt.Errorf("create history: %w", err)
		}

		result = mapStockOutOutput(s)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Cancel cancels a stock-out, releasing its reservation.
func (uc *StockOutUseCase) Cancel(ctx context.Context, id uuid.UUID) error {
	_, err := uc.Execute(ctx, id, valueobject.StockOutCancelled)
	return err
}

// List returns a paginated list of stock-outs.
func (uc *StockOutUseCase) List(ctx context.Context, filter port.StockOutFilter, page port.Pagination) (*dto.StockOutListOutput, error) {
	items, total, err := uc.stockOutRepo.FindAll(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("list stock outs: %w", err)
	}

	out := make([]*dto.StockOutOutput, len(items))
	for i, s := range items {
		out[i] = mapStockOutOutput(s)
	}
	return &dto.StockOutListOutput{Items: out, Total: total}, nil
}

func mapStockOutOutput(s *entity.StockOut) *dto.StockOutOutput {
	return &dto.StockOutOutput{
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

// GetDetail retrieves a detailed stock-out record including product info and timeline.
func (uc *StockOutUseCase) GetDetail(ctx context.Context, id uuid.UUID) (*dto.StockOutDetailOutput, error) {
	s, err := uc.stockOutRepo.FindByID(ctx, id)
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

	return &dto.StockOutDetailOutput{
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
		Timeline:    buildStockOutTimeline(s, histories),
	}, nil
}

func buildStockOutTimeline(s *entity.StockOut, histories []*entity.History) []dto.TimelineEntry {
	entries := []dto.TimelineEntry{{
		Status:      "DRAFT",
		Description: "Transaction record created by System API.",
		OccurredAt:  s.CreatedAt,
	}}
	for _, h := range histories {
		entries = append(entries, dto.TimelineEntry{
			Status:      h.NewStatus,
			Description: stockOutStatusDescription(h.NewStatus, h.OldStatus, s.PerformedBy),
			OccurredAt:  h.CreatedAt,
		})
	}
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries
}

func stockOutStatusDescription(newStatus, oldStatus, actor string) string {
	if actor == "" {
		actor = "System"
	}
	switch newStatus {
	case "IN_PROGRESS":
		return "Stock reservation confirmed and items being picked from warehouse."
	case "DONE":
		return fmt.Sprintf("Transaction completed and stock deducted from inventory by %s.", actor)
	case "CANCELLED":
		return fmt.Sprintf("Transaction cancelled from %s by %s.", oldStatus, actor)
	default:
		return fmt.Sprintf("Status changed from %s to %s.", oldStatus, newStatus)
	}
}
