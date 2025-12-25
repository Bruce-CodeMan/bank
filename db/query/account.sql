-- name: GetAccountById :one
SELECT *
FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountByUUID :one
SELECT *
FROM account
WHERE public_id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * 
FROM account
WHERE id = $1 LIMIT 1
FOR UPDATE;


-- name: ListAccounts :many
SELECT * FROM account 
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + @amount
WHERE id = @id
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account WHERE id = $1;

-- name: CreateAccount :one
INSERT INTO account (
    public_id,
    balance,
    currency,
    primary_user_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING
    id,
    public_id,
    balance,
    currency,
    created_at;