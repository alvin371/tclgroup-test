-- name: CreateStockIn :one
INSERT INTO stock_ins (id, product_id, quantity, status, notes, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetStockInByID :one
SELECT * FROM stock_ins WHERE id = $1 LIMIT 1;

-- name: UpdateStockIn :one
UPDATE stock_ins
SET status = $2, notes = $3, updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListStockIns :many
SELECT * FROM stock_ins
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountStockIns :one
SELECT COUNT(*) FROM stock_ins;

-- name: ListStockInsByStatus :many
SELECT * FROM stock_ins
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountStockInsByStatus :one
SELECT COUNT(*) FROM stock_ins WHERE status = $1;
