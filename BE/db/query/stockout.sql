-- name: CreateStockOut :one
INSERT INTO stock_outs (id, product_id, quantity, status, notes, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetStockOutByID :one
SELECT * FROM stock_outs WHERE id = $1 LIMIT 1;

-- name: UpdateStockOut :one
UPDATE stock_outs
SET status = $2, notes = $3, updated_at = $4
WHERE id = $1
RETURNING *;

-- name: ListStockOuts :many
SELECT * FROM stock_outs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountStockOuts :one
SELECT COUNT(*) FROM stock_outs;

-- name: ListStockOutsByStatus :many
SELECT * FROM stock_outs
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountStockOutsByStatus :one
SELECT COUNT(*) FROM stock_outs WHERE status = $1;
