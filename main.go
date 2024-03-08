package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/takeshiemoto/go-simple-api/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

// テスト容易性を上げるためにrun関数に切り出す
// context.Contextは複数の関数、ゴルーチン間でキャンセルシグナルを伝播させる手段を提供する
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)
	s := &http.Server{
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
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		os.Exit(1)
	}
}
