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
)

// ReservationRepo implements port.ReservationRepo.
type ReservationRepo struct {
	pool *pgxpool.Pool
}

// NewReservationRepo creates a new ReservationRepo.
func NewReservationRepo(pool *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{pool: pool}
}

// Create inserts a reservation row.
func (r *ReservationRepo) Create(ctx context.Context, res *entity.Reservation) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO reservations (id, stock_out_id, product_id, quantity, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		res.ID, res.StockOutID, res.ProductID, res.Quantity, res.Status, res.CreatedAt, res.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert reservation: %w", err)
	}
	return nil
}

// FindByStockOutID retrieves the reservation for a given stock-out.
func (r *ReservationRepo) FindByStockOutID(ctx context.Context, stockOutID uuid.UUID) (*entity.Reservation, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, stock_out_id, product_id, quantity, status, created_at, updated_at
		 FROM reservations WHERE stock_out_id = $1`, stockOutID)

	res := &entity.Reservation{}
	var status string
	err := row.Scan(&res.ID, &res.StockOutID, &res.ProductID, &res.Quantity, &status, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReservationNotFound
		}
		return nil, fmt.Errorf("scan reservation: %w", err)
	}
	res.Status = valueobject.ReservationStatus(status)
	return res, nil
}

// Update persists reservation status changes.
func (r *ReservationRepo) Update(ctx context.Context, res *entity.Reservation) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`UPDATE reservations SET status = $2, updated_at = $3 WHERE id = $1`,
		res.ID, res.Status, res.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update reservation: %w", err)
	}
	return nil
}
