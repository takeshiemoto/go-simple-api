package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	// httptest.NewRecorderは、HTTPレスポンスを受け取るための構造体を作成する。
	// テスト用のResponseWriterとして使用される。
	w := httptest.NewRecorder()
	// httptest.NewRequestは、テスト用のHTTPリクエストを作成する。
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	// テスト対象のMuxを作成する。
	sut := NewMux()
	// テスト対象のMuxのServeHTTPメソッドを呼び出し、
	// レスポンスをレコーダーに書き込む。
	sut.ServeHTTP(w, r)
	// レコーダーに書き込まれたレスポンスを取得する。
	resp := w.Result()
	// テスト終了時にレスポンスボディを閉じる。
	// テスト終了時にクリーンアップしないとリソースリークを引き起こす
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	// ステータスコードが200(OK)であることを確認する。
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status code 200, but", resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status" "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
