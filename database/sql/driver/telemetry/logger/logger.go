package logger

import (
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/ngrok/sqlmw"
)

// NewInterceptor for logger.
func NewInterceptor(name string, logger *logger.Logger, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{name: name, logger: logger, interceptor: interceptor}
}

// Interceptor for logger.
type Interceptor struct {
	interceptor sqlmw.Interceptor
	logger      *logger.Logger
	name        string
}

// ConnBeginTx for logger.
func (i *Interceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (context.Context, driver.Tx, error) {
	return i.interceptor.ConnBeginTx(ctx, conn, txOpts)
}

// ConnPrepareContext for logger.
func (i *Interceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (context.Context, driver.Stmt, error) {
	return i.interceptor.ConnPrepareContext(ctx, conn, query)
}

// ConnPing for logger.
func (i *Interceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	return i.interceptor.ConnPing(ctx, conn)
}

// ConnExecContext for logger.
func (i *Interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	attrs := []logger.Attr{
		logger.String(meta.SystemKey, i.name),
	}
	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))

	i.logger.Log(ctx, logger.NewMessage(message("exec conn"), err), attrs...)
	return res, err
}

// ConnQueryContext for logger.
func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now()
	attrs := []logger.Attr{
		logger.String(meta.SystemKey, i.name),
	}
	ctx, res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))

	i.logger.Log(ctx, logger.NewMessage(message("query conn"), err), attrs...)
	return ctx, res, err
}

// ConnectorConnect for logger.
func (i *Interceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
	return i.interceptor.ConnectorConnect(ctx, connect)
}

// ResultLastInsertId for logger.
func (i *Interceptor) ResultLastInsertId(res driver.Result) (int64, error) {
	return i.interceptor.ResultLastInsertId(res)
}

// ResultRowsAffected for logger.
func (i *Interceptor) ResultRowsAffected(res driver.Result) (int64, error) {
	return i.interceptor.ResultRowsAffected(res)
}

// RowsNext for logger.
func (i *Interceptor) RowsNext(ctx context.Context, rows driver.Rows, dest []driver.Value) error {
	return i.interceptor.RowsNext(ctx, rows, dest)
}

// RowsClose for logger.
func (i *Interceptor) RowsClose(ctx context.Context, rows driver.Rows) error {
	return i.interceptor.RowsClose(ctx, rows)
}

// StmtExecContext for logger.
func (i *Interceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	attrs := []logger.Attr{
		logger.String(meta.SystemKey, i.name),
	}
	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))

	i.logger.Log(ctx, logger.NewMessage(message("exec statement"), err), attrs...)
	return res, err
}

// StmtQueryContext for logger.
func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now()
	attrs := []logger.Attr{
		logger.String(meta.SystemKey, i.name),
	}
	ctx, res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)
	attrs = append(attrs, logger.String(meta.DurationKey, time.Since(start).String()))

	i.logger.Log(ctx, logger.NewMessage(message("query statement"), err), attrs...)
	return ctx, res, err
}

// StmtClose for logger.
func (i *Interceptor) StmtClose(ctx context.Context, stmt driver.Stmt) error {
	return i.interceptor.StmtClose(ctx, stmt)
}

// TxCommit for logger.
func (i *Interceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	return i.interceptor.TxCommit(ctx, tx)
}

// TxRollback for logger.
func (i *Interceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	return i.interceptor.TxRollback(ctx, tx)
}

func message(msg string) string {
	return "sql: " + msg
}
