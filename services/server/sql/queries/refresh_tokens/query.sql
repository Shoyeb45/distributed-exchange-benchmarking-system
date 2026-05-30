-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1 LIMIT 1;

-- name: DeleteRefreshTokenByHash :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;
