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

-- name: GetUserByUUID :one
SELECT users.*
FROM users 
WHERE users.id = $1
LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET email = $2, hashed_password = $3
WHERE users.id = $1;

-- name: SetUserRedState :exec
UPDATE users
SET is_chirpy_red = $2
WHERE users.id = $1;

-- name: GetUserRedState :one
SELECT users.is_chirpy_red
FROM users
WHERE users.id = $1;

-- name: ResetUsers :exec
TRUNCATE users CASCADE;

