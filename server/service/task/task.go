package task

import (
	"context"
	"fmt"

	db "github.com/punkzberryz/todo/db/sqlc"
)

var ErrOwnerNotMatched = fmt.Errorf("ownwer id does not match user id")

type Task struct {
	Store db.Store
}

// Create task from db
func (t *Task) CreateTask(ctx context.Context, arg db.CreateTaskParams) (*db.Task, error) {
	task, err := t.Store.CreateTask(ctx, arg)

	return &task, err
}

// Get task by Id
func (t *Task) GetTaskById(ctx context.Context, id int64, ownerId int64) (*db.Task, error) {
	task, err := t.Store.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	//check if task belongs to user
	if task.OwnerID != ownerId {
		return nil, ErrOwnerNotMatched
	}
	return &task, nil
}

// Get task list
func (t *Task) GetTaskList(ctx context.Context, ownerId int64, limit int32, pageId int32) ([]db.Task, error) {
	taskList, err := t.Store.GetTaskList(ctx, db.GetTaskListParams{
		OwnerID: ownerId,
		Limit:   limit,
		Offset:  (pageId - 1) * limit,
	})
	return taskList, err
}

// Update task by Id and OwnerId
func (t *Task) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (*db.Task, error) {
	task, err := t.Store.UpdateTask(ctx, arg)
	return &task, err
}

// Delete task by Id and OwnerId
func (t *Task) DeleteTask(ctx context.Context, arg db.DeleteTaskParams) error {
	err := t.Store.DeleteTask(ctx, arg)
	return err
}
