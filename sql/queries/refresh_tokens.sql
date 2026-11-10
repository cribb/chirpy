-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    NULL
)
RETURNING *;


-- name: GetRefreshTokenByUUID :one
SELECT refresh_tokens.*
FROM refresh_tokens
WHERE refresh_tokens.user_id = $1
LIMIT 1;


-- name: ResetRefreshTokens :exec
TRUNCATE refresh_tokens;



