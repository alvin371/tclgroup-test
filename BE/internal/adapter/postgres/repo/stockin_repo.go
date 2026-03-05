package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/domain/inventory/entity"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// StockInRepo implements port.StockInRepo.
type StockInRepo struct {
	pool *pgxpool.Pool
}

// NewStockInRepo creates a new StockInRepo.
func NewStockInRepo(pool *pgxpool.Pool) *StockInRepo {
	return &StockInRepo{pool: pool}
}

// Create inserts a stock_in row.
func (r *StockInRepo) Create(ctx context.Context, s *entity.StockIn) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO stock_ins (id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		s.ID, s.ProductID, s.Quantity, s.Status, s.Notes, s.UnitPrice, s.PerformedBy, s.Location, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert stock_in: %w", err)
	}
	return nil
}

// FindByID retrieves a stock-in by primary key.
func (r *StockInRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.StockIn, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at
		 FROM stock_ins WHERE id = $1`, id)

	s := &entity.StockIn{}
	var status string
	err := row.Scan(&s.ID, &s.ProductID, &s.Quantity, &status, &s.Notes, &s.UnitPrice, &s.PerformedBy, &s.Location, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrStockInNotFound
		}
		return nil, fmt.Errorf("scan stock_in: %w", err)
	}
	s.Status = valueobject.StockInStatus(status)
	return s, nil
}

// Update persists status and notes changes.
func (r *StockInRepo) Update(ctx context.Context, s *entity.StockIn) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`UPDATE stock_ins SET status = $2, notes = $3, updated_at = $4 WHERE id = $1`,
		s.ID, s.Status, s.Notes, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update stock_in: %w", err)
	}
	return nil
}

// FindAll returns a paginated list with optional filters.
func (r *StockInRepo) FindAll(ctx context.Context, filter port.StockInFilter, page port.Pagination) ([]*entity.StockIn, int64, error) {
	q := getQuerier(ctx, r.pool)

	args := []interface{}{}
	where := "WHERE 1=1"
	idx := 1

	if filter.ProductID != nil {
		where += fmt.Sprintf(" AND product_id = $%d", idx)
		args = append(args, *filter.ProductID)
		idx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, *filter.Status)
		idx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := q.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM stock_ins %s", where), countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count stock_ins: %w", err)
	}

	args = append(args, page.Limit, (page.Page-1)*page.Limit)
	rows, err := q.Query(ctx,
		fmt.Sprintf(`SELECT id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at
		             FROM stock_ins %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, idx, idx+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list stock_ins: %w", err)
	}
	defer rows.Close()

	var items []*entity.StockIn
	for rows.Next() {
		s := &entity.StockIn{}
		var status string
		if err = rows.Scan(&s.ID, &s.ProductID, &s.Quantity, &status, &s.Notes, &s.UnitPrice, &s.PerformedBy, &s.Location, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan stock_in row: %w", err)
		}
		s.Status = valueobject.StockInStatus(status)
		items = append(items, s)
	}
	return items, total, nil
}
