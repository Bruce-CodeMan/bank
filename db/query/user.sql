-- name: CreateUser :one
INSERT INTO "user" (
    public_id,
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    $1, $2, $3, $4, $5
) 
RETURNING 
    id,
    public_id,
    username,
    full_name,
    email;

-- name: GetUserByPublicID :one
SELECT * FROM "user"
WHERE public_id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;