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

-- name: GetUserByNameOrEmail :one
SELECT * FROM users
WHERE username = $1 OR email = $2
LIMIT 1;
