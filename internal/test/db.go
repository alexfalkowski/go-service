package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/driver"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/io"
)

// BenchmarkSQLDriver is a database/sql driver test double that returns empty result sets.
type BenchmarkSQLDriver struct{}

// Open implements driver.Driver.
func (BenchmarkSQLDriver) Open(string) (driver.Conn, error) {
	return BenchmarkSQLConn{}, nil
}

// BenchmarkSQLConn is a database/sql connection test double.
type BenchmarkSQLConn struct{}

// Prepare implements driver.Conn and returns driver.ErrSkip.
func (BenchmarkSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, driver.ErrSkip
}

// Close implements driver.Conn and always succeeds.
func (BenchmarkSQLConn) Close() error {
	return nil
}

// Begin implements driver.Conn and returns driver.ErrSkip.
func (BenchmarkSQLConn) Begin() (driver.Tx, error) {
	return nil, driver.ErrSkip
}

// QueryContext implements driver.QueryerContext and returns empty rows.
func (BenchmarkSQLConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return BenchmarkSQLRows{}, nil
}

// BenchmarkSQLRows is an empty rows test double.
type BenchmarkSQLRows struct{}

// Columns returns the value column.
func (BenchmarkSQLRows) Columns() []string {
	return []string{"value"}
}

// Close implements driver.Rows and always succeeds.
func (BenchmarkSQLRows) Close() error {
	return nil
}

// Next implements driver.Rows and returns io.EOF.
func (BenchmarkSQLRows) Next([]driver.Value) error {
	return io.EOF
}

func (w *World) registerDatabase() {
	pg.Register()

	w.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			if w.PG == nil || !w.PG.IsEnabled() {
				return nil
			}

			return w.openDatabase()
		},
		OnStop: func(_ context.Context) error {
			if w.DB == nil {
				return nil
			}

			return w.DB.Destroy()
		},
	})
}

func (w *World) openDatabase() error {
	db, err := pg.Connect(FS, w.PG)
	if err != nil {
		return err
	}

	w.DB = db

	return nil
}
