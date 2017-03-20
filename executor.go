package gorm

import (
	"database/sql"
	"errors"
	"golang.org/x/net/context"
)

// Executor wraps the underlying database connection, and adapts it
// to a context-aware interface.  The interface is based on go1.8's new context-aware
// methods, but backported to pre-go1.7 contexts.  This normalizes the interface to
// the database connection, enabling use of 1.7 and 1.8 features, if available.
//
// Note that go1.8 *sql.DB doesn't implement this.  *sql.DB has similar methods, but
// they use the "context" package.  This interface uses "golang.org/x/context" to maintain
// compatibility with pre go1.7 versions.  gorm accepts external connections which
// implement SQLCommon and may SQLCommonContext, then adapts them to this interface
// for internal use, and for use by Dialects.
type Executor interface {
	// ExecContext is a context-aware variant of Exec
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// PrepareContext is a context-aware variant of Prepare
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// QueryContext is a context-aware variant of Query
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// QueryRowContext is a context-aware variant of QueryRow
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// private

type executor struct {
	sqlCommon SQLCommon
}

func newExecutor(s SQLCommon) *executor {
	return &executor{sqlCommon: s}
}

type closer interface {
	Close() error
}

func (e *executor) Close() error {
	switch db := e.sqlCommon.(type) {
	case closer:
		return db.Close()
	default:
		return errors.New("can't close current db")
	}
}

type pinger interface {
	Ping() error
}

func (e *executor) ping() error {
	switch db := e.sqlCommon.(type) {
	case pinger:
		return db.Ping()
	default:
		return nil
	}
}

type sqlTx interface {
	SQLCommon
	Commit() error
	Rollback() error
}

func (e *executor) commit() error {
	switch db := e.sqlCommon.(type) {
	case sqlTx:
		return db.Commit()
	}
	return ErrInvalidTransaction
}

func (e *executor) rollback() error {
	switch db := e.sqlCommon.(type) {
	case sqlTx:
		return db.Rollback()
	default:
		return ErrInvalidTransaction
	}
}
