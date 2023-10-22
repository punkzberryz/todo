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

