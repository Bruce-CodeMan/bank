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
RETURNING *;

-- name: GetUserByPublicID :one
SELECT * FROM "user"
WHERE public_id = $1 LIMIT 1;