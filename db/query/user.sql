-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  user_role
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;
