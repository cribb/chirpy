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


-- -- GetChirpsForUser :many
-- SELECT chirps.*
-- FROM chirps
-- INNER JOIN users
--   ON posts.user_id = users.id
-- WHERE user_id = $1
-- LIMIT $2;

