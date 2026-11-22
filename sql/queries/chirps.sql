-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetChirps :exec
TRUNCATE TABLE chirps;


-- name: GetAllChirps :many
SELECT chirps.*
FROM chirps 
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT chirps.*
FROM chirps 
INNER JOIN users
  on chirps.user_id = users.id
WHERE chirps.id = $1
LIMIT 1;

-- name: GetChirpsForUser :many
SELECT chirps.*
FROM chirps
WHERE chirps.user_id = $1
ORDER BY created_at;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE chirps.id = $1
AND user_id = $2;

