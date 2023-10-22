-- name: CreateTask :one
INSERT INTO tasks (
    body,
    owner_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: GetTaskList :many
SELECT * FROM tasks
WHERE
    owner_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;