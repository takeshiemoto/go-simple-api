package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// ResponseJSON はJSON形式のレスポンスを生成する
func ResponseJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	// レスポンスのContent-typeを設定
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// bodyをJSON形式に変換
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := ErrorResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			fmt.Printf("write error response error: %v", err)
		}
		return
	}

	// ステータスコードを指定
	w.WriteHeader(status)
	// レスポンスボディをResponseWriterに書き込む
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
