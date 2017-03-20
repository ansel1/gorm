// +build go1.7

package gorm

import (
	"database/sql"
	"golang.org/x/net/context"
)

// ExecContext implements Executor
func (e *executor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch db := e.sqlCommon.(type) {
	case SQLCommonContext:
		return db.ExecContext(ctx, query, args...)
	default:
		return e.sqlCommon.Exec(query, args...)

	}
}

// PrepareContext implements Executor
func (e *executor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	switch db := e.sqlCommon.(type) {
	case SQLCommonContext:
		return db.PrepareContext(ctx, query)
	default:
		return e.sqlCommon.Prepare(query)

	}
}

// QueryContext implements Executor
func (e *executor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch db := e.sqlCommon.(type) {
	case SQLCommonContext:
		return db.QueryContext(ctx, query, args...)
	default:
		return e.sqlCommon.Query(query, args...)

	}
}

// QueryRowContext implements Executor
func (e *executor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch db := e.sqlCommon.(type) {
	case SQLCommonContext:
		return db.QueryRowContext(ctx, query, args...)
	default:
		return e.sqlCommon.QueryRow(query, args...)

	}
}
