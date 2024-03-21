package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

// テスト容易性を上げるためにrun関数に切り出す
// context.Contextは複数の関数、ゴルーチン間でキャンセルシグナルを伝播させる手段を提供する
func (s *Server) Run(ctx context.Context) error {
	// グレースフルシャットダウンの実装
	// プロセス終了シグナルを受け取ったときに実行中の処理を正しく終了させる。
	// シグナル(os.Interrupt, syscall.SIGTERM)を受け取ったらcontextをキャンセルする。
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// 指定されたフォーマットでログに出力する
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// グレースフルシャットダウンの終了を待つ
	return eg.Wait()
}
