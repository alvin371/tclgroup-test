package inventory

import (
	"context"
	"fmt"

	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// ReportUseCase handles reporting on completed stock transactions.
type ReportUseCase struct {
	stockInRepo  port.StockInRepo
	stockOutRepo port.StockOutRepo
}

// NewReportUseCase creates a new ReportUseCase.
func NewReportUseCase(stockInRepo port.StockInRepo, stockOutRepo port.StockOutRepo) *ReportUseCase {
	return &ReportUseCase{
		stockInRepo:  stockInRepo,
		stockOutRepo: stockOutRepo,
	}
}

// StockInReport returns only DONE stock-in transactions.
func (uc *ReportUseCase) StockInReport(ctx context.Context, filter port.StockInFilter, page port.Pagination) (*dto.StockInListOutput, error) {
	doneStatus := string(valueobject.StockInDone)
	filter.Status = &doneStatus

	items, total, err := uc.stockInRepo.FindAll(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("stock in report: %w", err)
	}

	out := make([]*dto.StockInOutput, len(items))
	for i, s := range items {
		out[i] = mapStockInOutput(s)
	}
	return &dto.StockInListOutput{Items: out, Total: total}, nil
}

// StockOutReport returns only DONE stock-out transactions.
func (uc *ReportUseCase) StockOutReport(ctx context.Context, filter port.StockOutFilter, page port.Pagination) (*dto.StockOutListOutput, error) {
	doneStatus := string(valueobject.StockOutDone)
	filter.Status = &doneStatus

	items, total, err := uc.stockOutRepo.FindAll(ctx, filter, page)
	if err != nil {
		return nil, fmt.Errorf("stock out report: %w", err)
	}

	out := make([]*dto.StockOutOutput, len(items))
	for i, s := range items {
		out[i] = mapStockOutOutput(s)
	}
	return &dto.StockOutListOutput{Items: out, Total: total}, nil
}
