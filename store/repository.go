package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/takeshiemoto/go-simple-api/clock"
	"github.com/takeshiemoto/go-simple-api/config"
)

// New はデータベースの接続を行う関数
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// fmt.Sprintfは指定されたフォーマットで文字列を生成する
	// parseTimeはtime.Time型のフィールドに正しい値をセットするために必要
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		return nil, nil, err
	}

	// コンテキストにタイムアウトを設定する
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// PingContextはデータベースへの接続を確認する
	if err := db.PingContext(ctx); err != nil {
		return nil, func() {
			// Pingに失敗した場合はデータベースの接続を閉じる
			_ = db.Close()
		}, err
	}
	// 接続を作成する
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

// Beginner はトランザクションを開始するインターフェース
type Beginner interface {
	// BeginTx はトランザクションを開始するために使われる
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Preparer はSQLステートメントを準備するインターフェース
type Preparer interface {
	// PreparexContext はSQLステートメントを準備する
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// Execer はSQLステートメントを実行するインターフェース
type Execer interface {
	// ExecContext はSQLステートメントを実行する
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// NamedExecContext は名前付きSQLステートメントを実行する
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}
