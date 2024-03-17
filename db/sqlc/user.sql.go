// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  full_name,
  hashed_password
) VALUES (
  $1, $2, $3, $4
) RETURNING id, username, hashed_password, full_name, email, user_role, password_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.FullName,
		arg.HashedPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.UserRole,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getListUsers = `-- name: GetListUsers :many
SELECT id, username, full_name, user_role, email, password_changed_at, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type GetListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetListUsersRow struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	UserRole          string    `json:"user_role"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (q *Queries) GetListUsers(ctx context.Context, arg GetListUsersParams) ([]GetListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getListUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetListUsersRow{}
	for rows.Next() {
		var i GetListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.FullName,
			&i.UserRole,
			&i.Email,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, username, hashed_password, full_name, email, user_role, password_changed_at, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.UserRole,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByNameOrEmail = `-- name: GetUserByNameOrEmail :one
SELECT id, username, hashed_password, full_name, email, user_role, password_changed_at, created_at FROM users
WHERE username = $1 OR email = $2
LIMIT 1
`

type GetUserByNameOrEmailParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) GetUserByNameOrEmail(ctx context.Context, arg GetUserByNameOrEmailParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByNameOrEmail, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.UserRole,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $1, password_changed_at = $2
WHERE id = $3
RETURNING id, username, hashed_password, full_name, email, user_role, password_changed_at, created_at
`

type UpdateUserPasswordParams struct {
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	ID                int64     `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.HashedPassword, arg.PasswordChangedAt, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.UserRole,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
