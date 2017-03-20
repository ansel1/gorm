// +build go1.8

package gorm_test

import (
	"context"
	"database/sql"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	"testing"
)

func (m *mockSQLCommonContext) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	m.calls = append(m.calls, acall{ctx: ctx, f: "begincontext"})
	return nil, errors.New("mock")
}

func (m *mockSQLCommonContext) Begin() (*sql.Tx, error) {
	m.calls = append(m.calls, acall{f: "begin"})
	return nil, errors.New("mock")
}

func TestBeginTx(t *testing.T) {
	sqlCommon := DB.CommonDB()
	mocker := &mockSQLCommonContext{inner: sqlCommon}
	db, err := gorm.Open(DB.Dialect().GetName(), mocker)
	if err != nil {
		t.Fatal("failed to open database connection")
	}
	ctx := context.WithValue(context.Background(), "color", "red")
	db2 := db.WithContext(ctx)
	db2.Begin()
	for _, c := range mocker.calls {
		if c.ctx == nil {
			t.Errorf("expected context to be passed to sql executor method: %s", c.f)
			break
		}
	}
	for _, n := range []string{"begincontext"} {
		var found bool
		for _, c := range mocker.calls {
			if c.f == n {
				found = true
			}
		}
		if !found {
			t.Errorf("expected %s to be called", n)
		}
	}
	for _, n := range []string{"begin"} {
		var found bool
		for _, c := range mocker.calls {
			if c.f == n {
				found = true
			}
		}
		if found {
			t.Errorf("expected %s to not be called", n)
		}
	}
}
