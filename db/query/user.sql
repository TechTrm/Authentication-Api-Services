-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  full_name,
  hashed_password
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetListUsers :many
SELECT id, username, full_name, user_role, email, password_changed_at, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetUserByNameOrEmail :one
SELECT * FROM users
WHERE username = $1 OR email = $2
LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $1, password_changed_at = $2
WHERE id = $3
RETURNING *;
