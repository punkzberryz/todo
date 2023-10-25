-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    email
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE
    id = $1 OR
    email = $2
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    email=sqlc.arg(new_email),
    username=$3,
    hashed_password=$4,
    password_changed_at=$5
WHERE id = $1 OR email = $2
RETURNING *;

