// +build go1.8

package gorm

import (
	"context"
	"database/sql"
)

type SQLBeginTxerContext interface {
	BeginTx(ctx context.Context, options *sql.TxOptions) (SQLTx, error)
}
