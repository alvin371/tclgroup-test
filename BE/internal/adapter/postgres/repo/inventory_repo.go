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
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// InventoryRepo implements port.InventoryRepo.
type InventoryRepo struct {
	pool *pgxpool.Pool
}

// NewInventoryRepo creates a new InventoryRepo.
func NewInventoryRepo(pool *pgxpool.Pool) *InventoryRepo {
	return &InventoryRepo{pool: pool}
}

// Create inserts an inventory row.
func (r *InventoryRepo) Create(ctx context.Context, inv *entity.Inventory) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO inventories (id, product_id, physical_stock, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		inv.ID, inv.ProductID, inv.PhysicalStock, inv.CreatedAt, inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert inventory: %w", err)
	}
	return nil
}

// FindByProductID retrieves the inventory for a product.
func (r *InventoryRepo) FindByProductID(ctx context.Context, productID uuid.UUID) (*entity.Inventory, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, product_id, physical_stock, created_at, updated_at
		 FROM inventories WHERE product_id = $1`, productID)

	inv := &entity.Inventory{}
	err := row.Scan(&inv.ID, &inv.ProductID, &inv.PhysicalStock, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInventoryNotFound
		}
		return nil, fmt.Errorf("scan inventory: %w", err)
	}
	return inv, nil
}

// LockByProductID retrieves the inventory with a SELECT FOR UPDATE NOWAIT lock.
func (r *InventoryRepo) LockByProductID(ctx context.Context, productID uuid.UUID) (*entity.Inventory, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, product_id, physical_stock, created_at, updated_at
		 FROM inventories WHERE product_id = $1 FOR UPDATE NOWAIT`, productID)

	inv := &entity.Inventory{}
	err := row.Scan(&inv.ID, &inv.ProductID, &inv.PhysicalStock, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInventoryNotFound
		}
		return nil, fmt.Errorf("lock inventory: %w", err)
	}
	return inv, nil
}

// Update persists changes to an inventory record.
func (r *InventoryRepo) Update(ctx context.Context, inv *entity.Inventory) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`UPDATE inventories SET physical_stock = $2, updated_at = $3 WHERE id = $1`,
		inv.ID, inv.PhysicalStock, inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}
	return nil
}

// GetTotalReserved sums the quantity of all ACTIVE reservations for a product.
func (r *InventoryRepo) GetTotalReserved(ctx context.Context, productID uuid.UUID) (int64, error) {
	q := getQuerier(ctx, r.pool)
	var total int64
	err := q.QueryRow(ctx,
		`SELECT COALESCE(SUM(quantity), 0) FROM reservations
		 WHERE product_id = $1 AND status = 'ACTIVE'`, productID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total reserved: %w", err)
	}
	return total, nil
}

// FindAllWithFilter returns a paginated list of inventory rows with product info.
func (r *InventoryRepo) FindAllWithFilter(ctx context.Context, filter port.InventoryFilter, page port.Pagination) ([]*port.InventoryRow, int64, error) {
	q := getQuerier(ctx, r.pool)

	args := []interface{}{}
	where := "WHERE 1=1"
	idx := 1

	if filter.ProductID != nil {
		where += fmt.Sprintf(" AND i.product_id = $%d", idx)
		args = append(args, *filter.ProductID)
		idx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := q.QueryRow(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM inventories i JOIN products p ON p.id = i.product_id %s", where),
		countArgs...,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count inventories: %w", err)
	}

	args = append(args, page.Limit, (page.Page-1)*page.Limit)
	rows, err := q.Query(ctx, fmt.Sprintf(`
		SELECT
			i.id, i.product_id, i.physical_stock, i.created_at, i.updated_at,
			p.id, p.sku, p.name, p.customer_id, p.created_at, p.updated_at,
			COALESCE((
				SELECT SUM(rv.quantity)
				FROM reservations rv
				WHERE rv.product_id = i.product_id AND rv.status = 'ACTIVE'
			), 0)::BIGINT AS reserved
		FROM inventories i
		JOIN products p ON p.id = i.product_id
		%s
		ORDER BY i.created_at DESC
		LIMIT $%d OFFSET $%d`, where, idx, idx+1), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list inventories: %w", err)
	}
	defer rows.Close()

	var result []*port.InventoryRow
	for rows.Next() {
		inv := &entity.Inventory{}
		prod := &entity.Product{}
		var reserved int64
		if err = rows.Scan(
			&inv.ID, &inv.ProductID, &inv.PhysicalStock, &inv.CreatedAt, &inv.UpdatedAt,
			&prod.ID, &prod.SKU, &prod.Name, &prod.CustomerID, &prod.CreatedAt, &prod.UpdatedAt,
			&reserved,
		); err != nil {
			return nil, 0, fmt.Errorf("scan inventory row: %w", err)
		}
		result = append(result, &port.InventoryRow{
			Inventory: inv,
			Product:   prod,
			Reserved:  reserved,
		})
	}
	return result, total, nil
}
