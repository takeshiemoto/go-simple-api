package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/takeshiemoto/go-simple-api/config"
)

// テスト容易性を上げるためにrun関数に切り出す
// context.Contextは複数の関数、ゴルーチン間でキャンセルシグナルを伝播させる手段を提供する
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// fmt.Sprintfは指定されたフォーマットで文字列を生成する
	// 文字列結合を行うケースで利用される
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		// エラーメッセージを出力しプログラムを終了する
		// %dは整数、%vは任意の型
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %v", url)

	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}
