package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/takeshiemoto/go-simple-api/entity"
	"github.com/takeshiemoto/go-simple-api/store"
)

type AddTask struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストからコンテキストを取得
	ctx := r.Context()
	// リクエストボディから取得するデータの構造体
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	// リクエストボディをでコードする
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		ResponseJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// データ構造のバリデーション
	err := at.Validator.Struct(b)
	if err != nil {
		ResponseJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// 新しいタスクの作成
	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}

	// ストアにタスクを追加する
	err = at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		ResponseJSON(ctx, w, &ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 追加したタスクのIDをレスポンスとして返す
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}

	ResponseJSON(ctx, w, rsp, http.StatusOK)
}
