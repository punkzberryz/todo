-- name: CreatePasswordResetSession :one
INSERT INTO password_reset_sessions (
    email,
    otp,
    expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPasswordResetSession :one
SELECT * FROM password_reset_sessions
WHERE  
    email =$1
LIMIT 1;

-- name: UpdatePasswordResetSession :one
UPDATE password_reset_sessions
SET
    otp = $2,
    expires_at = $3
WHERE
    email = $1
RETURNING *;

-- name: DeletePasswordResetSession :exec
DELETE FROM password_reset_sessions
WHERE email=$1;

