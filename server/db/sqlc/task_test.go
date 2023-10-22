package db

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/punkzberryz/todo/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTask(t *testing.T, user User) Task {
	arg := CreateTaskParams{
		OwnerID: user.ID,
		Body:    util.RandomString(10),
	}
	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)

	require.Equal(t, user.ID, task.OwnerID)
	require.Equal(t, arg.Body, task.Body)
	require.Equal(t, false, task.IsDone)
	require.NotZero(t, task.CreatedAt)

	return task
}

func TestCreateTask(t *testing.T) {
	user := CreateRandomUser(t)
	CreateRandomTask(t, user)
}

func TestGetTaskById(t *testing.T) {
	user := CreateRandomUser(t)
	task1 := CreateRandomTask(t, user)
	task2, err := testQueries.GetTask(context.Background(), task1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, task2)

	require.Equal(t, task1.Body, task2.Body)
	require.Equal(t, task1.IsDone, task2.IsDone)
	require.Equal(t, task1.OwnerID, task2.OwnerID)
	require.WithinDuration(t, task1.CreatedAt, task2.CreatedAt, time.Second)
}

func TestGetTaskList(t *testing.T) {
	user := CreateRandomUser(t)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			CreateRandomTask(t, user)
			wg.Done()
		}()
	}
	wg.Wait()

	taskList, err := testQueries.GetTaskList(context.Background(), GetTaskListParams{
		OwnerID: user.ID,
		Limit:   10,
		Offset:  1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, taskList)

	wg.Add(len(taskList))
	for _, task := range taskList {
		go func(task Task) {
			require.NotEmpty(t, task)
			require.Equal(t, user.ID, task.OwnerID)
			wg.Done()
		}(task)
	}
	wg.Wait()
}
