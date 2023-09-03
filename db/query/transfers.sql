-- name: NewTransfer :one
INSERT INTO transfers(
                      amount,
                      sender_id,
                      receiver_id
)VALUES ($1,$2,$3)RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers WHERE id=$1 LIMIT 1;

-- name: UpdateTransfer :one
UPDATE transfers SET amount=$2 WHERE id=$1 RETURNING *;

-- name: GetTransfers :many
SELECT * FROM transfers ORDER BY id LIMIT $1 OFFSET $2;

-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id=$1;