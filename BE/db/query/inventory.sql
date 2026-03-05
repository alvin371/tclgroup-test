-- name: CreateInventory :one
INSERT INTO inventories (id, product_id, physical_stock, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInventoryByProductID :one
SELECT * FROM inventories WHERE product_id = $1 LIMIT 1;

-- name: LockInventoryByProductID :one
SELECT * FROM inventories WHERE product_id = $1 LIMIT 1 FOR UPDATE NOWAIT;

-- name: UpdateInventory :one
UPDATE inventories
SET physical_stock = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: GetTotalReserved :one
SELECT COALESCE(SUM(quantity), 0)::BIGINT
FROM reservations
WHERE product_id = $1 AND status = 'ACTIVE';

-- name: ListInventoriesWithProduct :many
SELECT
    i.*,
    p.sku  AS product_sku,
    p.name AS product_name,
    COALESCE((
        SELECT SUM(r.quantity)
        FROM reservations r
        WHERE r.product_id = i.product_id AND r.status = 'ACTIVE'
    ), 0)::BIGINT AS reserved
FROM inventories i
JOIN products p ON p.id = i.product_id
ORDER BY i.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountInventories :one
SELECT COUNT(*) FROM inventories;
