package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/takeshiemoto/go-simple-api/entity"
	"github.com/takeshiemoto/go-simple-api/store"
)

type ListTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	if err != nil {
		ResponseJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	var rsp []task
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	ResponseJSON(ctx, w, rsp, http.StatusOK)
}
