package gorm

import (
	"database/sql"
)

// SQLCommon is the minimal database connection functionality gorm requires.  Implemented by *sql.DB.
type SQLCommon interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// RawDBer is an interface for things which can return a raw *sql.DB connection.
// It is intended to be used when callers pass a custom SQLCommon or SQLCommonContext
// implementation, but still want gorm.DB.DB() to return an underlying *sql.DB.
type RawDBer interface {
	RawDB() *sql.DB
}
