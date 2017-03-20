// +build !go1.7

package gorm

import (
	"database/sql"
	"golang.org/x/net/context"
)

// ExecContext implements Executor
func (e *executor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return e.sqlCommon.Exec(query, args...)
}

// PrepareContext implements Executor
func (e *executor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return e.sqlCommon.Prepare(query)
}

// QueryContext implements Executor
func (e *executor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return e.sqlCommon.Query(query, args...)
}

// QueryRowContext implements Executor
func (e *executor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return e.sqlCommon.QueryRow(query, args...)
}

type sqlDb interface {
	Begin() (*sql.Tx, error)
}

func (e *executor) begin(ctx context.Context) (*executor, error) {
	switch db := e.sqlCommon.(type) {
	case sqlDb:
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		return newExecutor(tx), nil
	default:
		return nil, ErrCantStartTransaction
	}
}
