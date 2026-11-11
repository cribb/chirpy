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


-- name: GetUserbyRefreshToken :one
SELECT users.*
FROM refresh_tokens
INNER JOIN users
  ON refresh_tokens.user_id = users.id
WHERE (refresh_tokens.token = $1
  AND refresh_tokens.expires_at > NOW()
  AND refresh_tokens.revoked_at IS NULL)
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE refresh_tokens.token = $1;


-- name: ResetRefreshTokens :exec
TRUNCATE refresh_tokens;

