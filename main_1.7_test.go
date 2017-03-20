// +build go1.7

package gorm_test

import (
	"database/sql"
	"testing"
	//gocontext "golang.org/x/net/context"
	"context"
	"github.com/jinzhu/gorm"
)

type acall struct {
	f   string
	ctx context.Context
}

type mockSQLCommonContext struct {
	inner gorm.SQLCommon
	calls []acall
}

func (m *mockSQLCommonContext) acall(name string, ctx context.Context) {
	m.calls = append(m.calls, acall{f: name, ctx: ctx})
}

func (m *mockSQLCommonContext) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.acall("execcontext", ctx)
	return m.inner.Exec(query, args...)
}

func (m *mockSQLCommonContext) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// not tested, not used by gorm
	return m.inner.Prepare(query)
}

func (m *mockSQLCommonContext) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.acall("querycontext", ctx)
	return m.inner.Query(query, args...)
}

func (m *mockSQLCommonContext) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.acall("queryrowcontext", ctx)
	return m.inner.QueryRow(query, args...)
}

func (m *mockSQLCommonContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.acall("exec", nil)
	return m.inner.Exec(query, args...)
}

func (m *mockSQLCommonContext) Prepare(query string) (*sql.Stmt, error) {
	return m.inner.Prepare(query)
}

func (m *mockSQLCommonContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	m.acall("query", nil)
	return m.inner.Query(query, args...)
}

func (m *mockSQLCommonContext) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.inner.QueryRow(query, args...)
}

func (m *mockSQLCommonContext) reset() {
	m.calls = nil
}

func TestSQLCommonContext(t *testing.T) {
	sqlCommon := DB.CommonDB()
	mocker := &mockSQLCommonContext{inner: sqlCommon}
	db, err := gorm.Open(DB.Dialect().GetName(), mocker)
	if err != nil {
		t.Fatal("failed to open database connection")
	}
	ctx := context.WithValue(context.Background(), "color", "red")
	db2 := db.WithContext(ctx)
	user1 := User{Name: "RowUser1", Age: 1, Birthday: parseTime("2000-1-1")}
	db2.Exec("select * from users;")
	db2.Save(user1)
	db2.First(&User{})
	db2.Model(&User{}).Row()
	for _, c := range mocker.calls {
		if c.ctx == nil {
			t.Errorf("expected context to be passed to sql executor method: %s", c.f)
			break
		}
	}
	for _, n := range []string{"execcontext", "querycontext", "queryrowcontext"} {
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
	for _, n := range []string{"exec", "query", "queryrow"} {
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
