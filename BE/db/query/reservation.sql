-- name: CreateReservation :one
INSERT INTO reservations (id, stock_out_id, product_id, quantity, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetReservationByStockOutID :one
SELECT * FROM reservations WHERE stock_out_id = $1 LIMIT 1;

-- name: UpdateReservation :one
UPDATE reservations
SET status = $2, updated_at = $3
WHERE id = $1
RETURNING *;
