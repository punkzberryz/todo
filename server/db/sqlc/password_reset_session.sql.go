// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: password_reset_session.sql

package db

import (
	"context"
	"time"
)

const createPasswordResetSession = `-- name: CreatePasswordResetSession :one
INSERT INTO password_reset_sessions (
    email,
    otp,
    expires_at
) VALUES (
    $1, $2, $3
) RETURNING email, otp, expires_at, created_at
`

type CreatePasswordResetSessionParams struct {
	Email     string    `json:"email"`
	Otp       string    `json:"otp"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (q *Queries) CreatePasswordResetSession(ctx context.Context, arg CreatePasswordResetSessionParams) (PasswordResetSession, error) {
	row := q.queryRow(ctx, q.createPasswordResetSessionStmt, createPasswordResetSession, arg.Email, arg.Otp, arg.ExpiresAt)
	var i PasswordResetSession
	err := row.Scan(
		&i.Email,
		&i.Otp,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deletePasswordResetSession = `-- name: DeletePasswordResetSession :exec
DELETE FROM password_reset_sessions
WHERE email=$1
`

func (q *Queries) DeletePasswordResetSession(ctx context.Context, email string) error {
	_, err := q.exec(ctx, q.deletePasswordResetSessionStmt, deletePasswordResetSession, email)
	return err
}

const getPasswordResetSession = `-- name: GetPasswordResetSession :one
SELECT email, otp, expires_at, created_at FROM password_reset_sessions
WHERE  
    email =$1
LIMIT 1
`

func (q *Queries) GetPasswordResetSession(ctx context.Context, email string) (PasswordResetSession, error) {
	row := q.queryRow(ctx, q.getPasswordResetSessionStmt, getPasswordResetSession, email)
	var i PasswordResetSession
	err := row.Scan(
		&i.Email,
		&i.Otp,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const updatePasswordResetSession = `-- name: UpdatePasswordResetSession :one
UPDATE password_reset_sessions
SET
    otp = $2,
    expires_at = $3
WHERE
    email = $1
RETURNING email, otp, expires_at, created_at
`

type UpdatePasswordResetSessionParams struct {
	Email     string    `json:"email"`
	Otp       string    `json:"otp"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (q *Queries) UpdatePasswordResetSession(ctx context.Context, arg UpdatePasswordResetSessionParams) (PasswordResetSession, error) {
	row := q.queryRow(ctx, q.updatePasswordResetSessionStmt, updatePasswordResetSession, arg.Email, arg.Otp, arg.ExpiresAt)
	var i PasswordResetSession
	err := row.Scan(
		&i.Email,
		&i.Otp,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
