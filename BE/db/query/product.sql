-- name: CreateProduct :one
INSERT INTO products (id, sku, name, customer_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;
