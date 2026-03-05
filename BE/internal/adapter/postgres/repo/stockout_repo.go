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

// StockOutRepo implements port.StockOutRepo.
type StockOutRepo struct {
	pool *pgxpool.Pool
}

// NewStockOutRepo creates a new StockOutRepo.
func NewStockOutRepo(pool *pgxpool.Pool) *StockOutRepo {
	return &StockOutRepo{pool: pool}
}

// Create inserts a stock_out row.
func (r *StockOutRepo) Create(ctx context.Context, s *entity.StockOut) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO stock_outs (id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		s.ID, s.ProductID, s.Quantity, s.Status, s.Notes, s.UnitPrice, s.PerformedBy, s.Location, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert stock_out: %w", err)
	}
	return nil
}

// FindByID retrieves a stock-out by primary key.
func (r *StockOutRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.StockOut, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at
		 FROM stock_outs WHERE id = $1`, id)

	s := &entity.StockOut{}
	var status string
	err := row.Scan(&s.ID, &s.ProductID, &s.Quantity, &status, &s.Notes, &s.UnitPrice, &s.PerformedBy, &s.Location, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrStockOutNotFound
		}
		return nil, fmt.Errorf("scan stock_out: %w", err)
	}
	s.Status = valueobject.StockOutStatus(status)
	return s, nil
}

// Update persists status and notes changes.
func (r *StockOutRepo) Update(ctx context.Context, s *entity.StockOut) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`UPDATE stock_outs SET status = $2, notes = $3, updated_at = $4 WHERE id = $1`,
		s.ID, s.Status, s.Notes, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update stock_out: %w", err)
	}
	return nil
}

// FindAll returns a paginated list with optional filters.
func (r *StockOutRepo) FindAll(ctx context.Context, filter port.StockOutFilter, page port.Pagination) ([]*entity.StockOut, int64, error) {
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
	if err := q.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM stock_outs %s", where), countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count stock_outs: %w", err)
	}

	args = append(args, page.Limit, (page.Page-1)*page.Limit)
	rows, err := q.Query(ctx,
		fmt.Sprintf(`SELECT id, product_id, quantity, status, notes, unit_price, performed_by, location, created_at, updated_at
		             FROM stock_outs %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, idx, idx+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list stock_outs: %w", err)
	}
	defer rows.Close()

	var items []*entity.StockOut
	for rows.Next() {
		s := &entity.StockOut{}
		var status string
		if err = rows.Scan(&s.ID, &s.ProductID, &s.Quantity, &status, &s.Notes, &s.UnitPrice, &s.PerformedBy, &s.Location, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan stock_out row: %w", err)
		}
		s.Status = valueobject.StockOutStatus(status)
		items = append(items, s)
	}
	return items, total, nil
}
