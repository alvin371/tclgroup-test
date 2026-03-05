-- name: CreateHistory :one
INSERT INTO histories (id, entity_type, entity_id, action, old_status, new_status, quantity, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetHistoriesByEntityID :many
SELECT * FROM histories
WHERE entity_id = $1
ORDER BY created_at ASC;
