package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tclgroup/stock-management/internal/domain/inventory/entity"
)

// HistoryRepo implements port.HistoryRepo.
type HistoryRepo struct {
	pool *pgxpool.Pool
}

// NewHistoryRepo creates a new HistoryRepo.
func NewHistoryRepo(pool *pgxpool.Pool) *HistoryRepo {
	return &HistoryRepo{pool: pool}
}

// Create inserts a history record (append-only).
func (r *HistoryRepo) Create(ctx context.Context, h *entity.History) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO histories (id, entity_type, entity_id, action, old_status, new_status, quantity, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		h.ID, h.EntityType, h.EntityID, h.Action, h.OldStatus, h.NewStatus, h.Quantity, h.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert history: %w", err)
	}
	return nil
}

// FindByEntityID retrieves all history records for an entity, ordered by time.
func (r *HistoryRepo) FindByEntityID(ctx context.Context, entityID uuid.UUID) ([]*entity.History, error) {
	q := getQuerier(ctx, r.pool)
	rows, err := q.Query(ctx,
		`SELECT id, entity_type, entity_id, action, old_status, new_status, quantity, created_at
		 FROM histories WHERE entity_id = $1 ORDER BY created_at ASC`, entityID)
	if err != nil {
		return nil, fmt.Errorf("query histories: %w", err)
	}
	defer rows.Close()

	var items []*entity.History
	for rows.Next() {
		h := &entity.History{}
		if err = rows.Scan(
			&h.ID, &h.EntityType, &h.EntityID, &h.Action,
			&h.OldStatus, &h.NewStatus, &h.Quantity, &h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan history row: %w", err)
		}
		items = append(items, h)
	}
	return items, nil
}
