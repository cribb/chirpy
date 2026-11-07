-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: GetUserByEmail :one
SELECT users.*
FROM users 
WHERE users.email = $1
LIMIT 1;


-- name: ResetUsers :exec
TRUNCATE users CASCADE;

