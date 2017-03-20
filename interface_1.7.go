// +build go1.7

package gorm

import (
	"context"
	"database/sql"
)

// SQLCommonContext is a context-aware variant of SQLCommon
type SQLCommonContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
type SQLTx interface {
	SQLCommon
	Commit() error
	Rollback() error
}

type SQLBeginner interface {
	Begin(ctx context.Context) (SQLTx, error)
}
