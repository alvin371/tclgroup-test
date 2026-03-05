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

// ProductRepo implements port.ProductRepo.
type ProductRepo struct {
	pool *pgxpool.Pool
}

// NewProductRepo creates a new ProductRepo.
func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

// Create inserts a product row.
func (r *ProductRepo) Create(ctx context.Context, p *entity.Product) error {
	q := getQuerier(ctx, r.pool)
	_, err := q.Exec(ctx,
		`INSERT INTO products (id, sku, name, customer_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID, p.SKU, p.Name, p.CustomerID, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

// FindByID retrieves a product by primary key.
func (r *ProductRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	q := getQuerier(ctx, r.pool)
	row := q.QueryRow(ctx,
		`SELECT id, sku, name, customer_id, created_at, updated_at
		 FROM products WHERE id = $1`, id)

	p := &entity.Product{}
	err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.CustomerID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("scan product: %w", err)
	}
	return p, nil
}

// FindAll returns a paginated list of products with optional filters.
func (r *ProductRepo) FindAll(ctx context.Context, filter port.ProductFilter, page port.Pagination) ([]*entity.Product, int64, error) {
	q := getQuerier(ctx, r.pool)

	// Build dynamic WHERE clause
	args := []interface{}{}
	where := "WHERE 1=1"
	idx := 1

	if filter.CustomerID != nil {
		where += fmt.Sprintf(" AND customer_id = $%d", idx)
		args = append(args, *filter.CustomerID)
		idx++
	}
	if filter.SKU != nil {
		where += fmt.Sprintf(" AND sku = $%d", idx)
		args = append(args, *filter.SKU)
		idx++
	}
	if filter.Name != nil {
		where += fmt.Sprintf(" AND name ILIKE $%d", idx)
		args = append(args, "%"+*filter.Name+"%")
		idx++
	}

	// Count
	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := q.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM products %s", where), countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	// List
	args = append(args, page.Limit, (page.Page-1)*page.Limit)
	rows, err := q.Query(ctx,
		fmt.Sprintf(`SELECT id, sku, name, customer_id, created_at, updated_at
		             FROM products %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, idx, idx+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		p := &entity.Product{}
		if err = rows.Scan(&p.ID, &p.SKU, &p.Name, &p.CustomerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, p)
	}
	return products, total, nil
}
