package main

import (
	"context"
	"fmt"
	"github.com/takeshiemoto/go-simple-api/config"
	"log"
	"net"
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

	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
