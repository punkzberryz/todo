package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/service/task"
	"github.com/punkzberryz/todo/service/token"
)

type TaskResponse struct {
	*db.Task
}

func (trsp *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getTaskIdFromURLPath(r *http.Request) (int, error) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		err := fmt.Errorf("task ID is required")
		return 0, err
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		err := fmt.Errorf("invalid task ID")
		return 0, err
	}
	return id, nil
}

// get task by id param
func (server *Server) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIdFromURLPath(r)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	payload := r.Context().Value(payloadKey).(*token.Payload)
	taskRsp, err := server.task.GetTaskById(r.Context(), int64(id), payload.User.ID)
	if err != nil {
		if err == task.ErrOwnerNotMatched {
			render.Render(w, r, ErrUnauthorized(err))
			return
		}
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, &TaskResponse{
		Task: taskRsp,
	}); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

type TaskListResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func (*TaskListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// get task list by owner id
// /task?pageId=1&limit=10
func (server *Server) getTaskList(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(payloadKey).(*token.Payload)

	queryStrings := r.URL.Query()
	pageId, err := strconv.Atoi(queryStrings.Get("pageId"))
	if err != nil {
		pageId = 1
	}
	limit, err := strconv.Atoi(queryStrings.Get("limit"))
	if err != nil {
		limit = 10
	}

	taskList, err := server.task.GetTaskList(r.Context(), payload.User.ID, int32(limit), int32(pageId))
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.Render(w, r, &TaskListResponse{Tasks: taskList}); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

type CreateTaskRequest struct {
	Body string `json:"body"`
}

// Create Bind function for Body request validation
func (c *CreateTaskRequest) Bind(r *http.Request) error {
	if c.Body == "" {
		return fmt.Errorf("body is a required field")
	}
	return nil
}

// create new task
func (server *Server) createTask(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(payloadKey).(*token.Payload)

	data := &CreateTaskRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	task, err := server.task.CreateTask(r.Context(),
		db.CreateTaskParams{
			OwnerID: payload.User.ID,
			Body:    data.Body,
		},
	)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.Render(w, r, &TaskResponse{Task: task}); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// update task
type UpdateTaskRequest struct {
	Body   string `json:"body"`
	IsDone bool   `json:"isDone"`
}

// Update Bind function for Body request validation
func (c *UpdateTaskRequest) Bind(r *http.Request) error {
	if c.Body == "" {
		return fmt.Errorf("body is a required field")
	}
	return nil
}

func (server *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIdFromURLPath(r)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	payload := r.Context().Value(payloadKey).(*token.Payload)
	data := &UpdateTaskRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	//We update task with Query with taskId and ownerId
	//if userId doesn't match ownerId, the query will fail,
	//hence we don't need to compare ownerId with userId anymore
	task, err := server.task.UpdateTask(r.Context(), db.UpdateTaskParams{
		ID:      int64(id),
		Body:    data.Body,
		IsDone:  data.IsDone,
		OwnerID: payload.User.ID,
	})
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	if err := render.Render(w, r, &TaskResponse{Task: task}); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// Delete task
// We improve response time by quering both
// taskId and userId,
// downside is that the response always return success
// even if taskId doesn't exist nor matching the userId
// (but the task won't be deleted if ownerId!=userId)
type deleteTaskResponse struct {
	Message string `json:"message"`
}

func (*deleteTaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (server *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getTaskIdFromURLPath(r)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	payload := r.Context().Value(payloadKey).(*token.Payload)
	err = server.task.DeleteTask(r.Context(),
		db.DeleteTaskParams{
			ID:      int64(id),
			OwnerID: payload.User.ID,
		})
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	rsp := &deleteTaskResponse{
		Message: fmt.Sprintf("delete task id %d success", id),
	}
	if err := render.Render(w, r, rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
