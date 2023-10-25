// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: task.sql

package db

import (
	"context"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (
    body,
    owner_id
) VALUES (
    $1, $2
) RETURNING id, body, is_done, owner_id, created_at
`

type CreateTaskParams struct {
	Body    string `json:"body"`
	OwnerID int64  `json:"ownerId"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.queryRow(ctx, q.createTaskStmt, createTask, arg.Body, arg.OwnerID)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.IsDone,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1 AND owner_id = $2
`

type DeleteTaskParams struct {
	ID      int64 `json:"id"`
	OwnerID int64 `json:"ownerId"`
}

func (q *Queries) DeleteTask(ctx context.Context, arg DeleteTaskParams) error {
	_, err := q.exec(ctx, q.deleteTaskStmt, deleteTask, arg.ID, arg.OwnerID)
	return err
}

const getTask = `-- name: GetTask :one
SELECT id, body, is_done, owner_id, created_at FROM tasks
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTask(ctx context.Context, id int64) (Task, error) {
	row := q.queryRow(ctx, q.getTaskStmt, getTask, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.IsDone,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}

const getTaskList = `-- name: GetTaskList :many
SELECT id, body, is_done, owner_id, created_at FROM tasks
WHERE
    owner_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type GetTaskListParams struct {
	OwnerID int64 `json:"ownerId"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) GetTaskList(ctx context.Context, arg GetTaskListParams) ([]Task, error) {
	rows, err := q.query(ctx, q.getTaskListStmt, getTaskList, arg.OwnerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Task{}
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Body,
			&i.IsDone,
			&i.OwnerID,
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

const updateTask = `-- name: UpdateTask :one
UPDATE tasks
SET 
    body = $3,
    is_done = $4
WHERE id = $1 AND owner_id = $2
RETURNING id, body, is_done, owner_id, created_at
`

type UpdateTaskParams struct {
	ID      int64  `json:"id"`
	OwnerID int64  `json:"ownerId"`
	Body    string `json:"body"`
	IsDone  bool   `json:"isDone"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.queryRow(ctx, q.updateTaskStmt, updateTask,
		arg.ID,
		arg.OwnerID,
		arg.Body,
		arg.IsDone,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.IsDone,
		&i.OwnerID,
		&i.CreatedAt,
	)
	return i, err
}
