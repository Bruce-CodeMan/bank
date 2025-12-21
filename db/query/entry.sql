-- name: CreateEntry :one
INSERT INTO entry (
    account_id,
    transfer_id,
    amount
) VALUES (
    $1, $2, $3
)
RETURNING
    id,
    account_id,
    transfer_id,
    amount,
    created_at;

-- name: GetEntry :one
SELECT *
FROM entry
WHERE id = $1;