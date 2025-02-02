-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, expires_on, created_on, updated_on)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetRefreshTokenInfoByToken :one
SELECT refresh_tokens.id, refresh_tokens.user_id, refresh_tokens.token, refresh_tokens.expires_on
  FROM refresh_tokens
  WHERE refresh_tokens.token = $1;

-- name: DeleteRefreshTokenById :exec
DELETE FROM refresh_tokens
WHERE refresh_tokens.id = $1;

