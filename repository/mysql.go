package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Result = sql.Result
type Row = sql.Row
type Rows = sql.Rows

type DBTX interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Exec(query string, args ...any) (Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (Result, error)
	Query(query string, args ...any) (*Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
	QueryRow(query string, args ...any) *Row
	QueryRowContext(ctx context.Context, query string, args ...any) *Row
}

func NewMysqlDB(connUrl string) *sql.DB {
	db, err := sql.Open("mysql", connUrl)
	if err != nil {
		log.Fatalf("[init] -- (repository/mysql) status: fail")
	} else {
		log.Printf("[init] -- (repository/mysql) status: success")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	if err := db.Ping(); err != nil {
		log.Fatalf("[conn] -- (repository/mysql) stauts: fail")
	} else {
		log.Printf("[conn] -- (repository/mysql) status: success")
	}
	return db
}
