package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/token"
)

type TaskResponse struct {
	*db.Task
}

func (trsp *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// get task by id param
func (server *Server) getTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		render.Render(w, r, ErrRender(fmt.Errorf("task ID is required")))
		return
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid task ID")))
		return
	}
	task, err := server.store.GetTask(r.Context(), int64(id))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, &TaskResponse{
		Task: &task,
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

	taskList, err := server.store.GetTaskList(r.Context(), db.GetTaskListParams{
		OwnerID: payload.User.ID,
		Limit:   int32(limit),
		Offset:  (int32(pageId) - 1) * int32(limit),
	})
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

	arg := db.CreateTaskParams{
		OwnerID: payload.User.ID,
		Body:    data.Body,
	}

	task, err := server.store.CreateTask(r.Context(), arg)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	if err := render.Render(w, r, &TaskResponse{Task: &task}); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
