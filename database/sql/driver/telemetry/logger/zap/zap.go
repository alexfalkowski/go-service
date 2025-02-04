package zap

import (
	"context"
	"database/sql/driver"
	"time"

	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/ngrok/sqlmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewInterceptor for zap.
func NewInterceptor(name string, logger *zap.Logger, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{name: name, logger: logger, interceptor: interceptor}
}

// Interceptor for zap.
type Interceptor struct {
	interceptor sqlmw.Interceptor
	logger      *zap.Logger
	name        string
}

func (i *Interceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (context.Context, driver.Tx, error) {
	return i.interceptor.ConnBeginTx(ctx, conn, txOpts)
}

func (i *Interceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (context.Context, driver.Stmt, error) {
	return i.interceptor.ConnPrepareContext(ctx, conn, query)
}

func (i *Interceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	return i.interceptor.ConnPing(ctx, conn)
}

func (i *Interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, i.name),
	}

	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("exec conn"), err, i.logger, fields...)

	return res, err
}

func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, i.name),
	}

	ctx, res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("query conn"), err, i.logger, fields...)

	return ctx, res, err
}

func (i *Interceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
	return i.interceptor.ConnectorConnect(ctx, connect)
}

//nolint:revive,stylecheck
func (i *Interceptor) ResultLastInsertId(res driver.Result) (int64, error) {
	return i.interceptor.ResultLastInsertId(res)
}

func (i *Interceptor) ResultRowsAffected(res driver.Result) (int64, error) {
	return i.interceptor.ResultRowsAffected(res)
}

func (i *Interceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) error {
	return i.interceptor.RowsNext(ctx, rows, dest)
}

func (i *Interceptor) RowsClose(ctx context.Context, rows driver.Rows) error {
	return i.interceptor.RowsClose(ctx, rows)
}

func (i *Interceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, i.name),
	}

	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)

	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("exec statement"), err, i.logger, fields...)

	return res, err
}

func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now()
	fields := []zapcore.Field{
		zap.String(meta.ServiceKey, i.name),
	}

	ctx, res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)
	fields = append(fields, tz.Meta(ctx)...)
	fields = append(fields, zap.Stringer(meta.DurationKey, time.Since(start)))

	tz.LogWithLogger(message("query statement"), err, i.logger, fields...)

	return ctx, res, err
}

func (i *Interceptor) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	return i.interceptor.StmtClose(ctx, stmt)
}

func (i *Interceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	return i.interceptor.TxCommit(ctx, tx)
}

func (i *Interceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	return i.interceptor.TxRollback(ctx, tx)
}

func message(msg string) string {
	return "sql: " + msg
}
