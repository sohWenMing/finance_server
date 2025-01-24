-- name: GetUser :many
SELECT * from users; 

-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
SELECT * from users
WHERE id = $1;


-- name: DeleteUserById :exec
DELETE from users
WHERE id = $1;