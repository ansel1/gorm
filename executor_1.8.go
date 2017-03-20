// +build go1.8

package gorm

import (
	"context"
	"database/sql"
)

type sqlDb interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (e *executor) begin(ctx context.Context) (*executor, error) {
	switch db := e.sqlCommon.(type) {
	case SQLBeginTxerContext:
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		return newExecutor(tx), nil
	case SQLBeginner:
		tx, err := db.Begin(ctx)
		if err != nil {
			return nil, err
		}
		return newExecutor(tx), nil
	case sqlDb:
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		return newExecutor(tx), nil
	default:
		return nil, ErrCantStartTransaction
	}
}
