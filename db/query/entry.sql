-- name: CreateEntry :one
INSERT INTO entries(
                    amount,
                    account_id
)VALUES ($1,$2)RETURNING *;

-- name: UpdateEntry :one
UPDATE entries SET amount=$2 WHERE id=$1 RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id=$1;

-- name: GetEntries :many
SELECT * FROM entries ORDER BY id LIMIT $1 OFFSET $2   ;

-- name: GetEntry :one
SELECT * FROM entries WHERE id=$1 LIMIT 1;