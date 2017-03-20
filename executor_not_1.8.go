// +build go1.7,!go1.8

package gorm

import (
	"database/sql"
	"golang.org/x/net/context"
)

type sqlDb interface {
	Begin() (*sql.Tx, error)
}

func (e *executor) begin(ctx context.Context) (*executor, error) {
	switch db := e.sqlCommon.(type) {
	case SQLBeginner:
		tx, err := db.Begin(ctx)
		if err != nil {
			return nil, err
		}
		return newExecutor(tx), nil
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
