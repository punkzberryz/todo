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
    username = COALESCE(sqlc.narg(username), username),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    email= COALESCE(sqlc.narg(new_email), email)           
WHERE id = $1 OR email = $2
RETURNING *;

