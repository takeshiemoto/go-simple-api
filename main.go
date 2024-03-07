package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

// テスト容易性を上げるためにrun関数に切り出す
// context.Contextは複数の関数、ゴルーチン間でキャンセルシグナルを伝播させる手段を提供する
func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/**
			指定されたio.Writerに実装されている出力ストリームに書き込む
			```go
			// 標準出力を例として使用
			w := os.Stdout
			fmt.Fprintf(w, "Hello, %s!", "world")
			*/
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// 指定されたフォーマットでログに出力する
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
