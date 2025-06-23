package tracer

import (
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/ngrok/sqlmw"
)

// NewInterceptor for tracer.
func NewInterceptor(driver string, tracer *tracer.Tracer, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{driver: driver, tracer: tracer, interceptor: interceptor}
}

// Interceptor for tracer.
type Interceptor struct {
	tracer      *tracer.Tracer
	interceptor sqlmw.Interceptor
	driver      string
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
	ctx, span := i.tracer.StartClient(ctx, operationName("exec conn"), attributes.String("db.sql.query", query))
	defer span.End()

	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return res, err
}

func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	ctx, span := i.tracer.StartClient(ctx, operationName("query conn"), attributes.String("db.sql.query", query))
	defer span.End()

	ctx, res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return ctx, res, err
}

func (i *Interceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
	return i.interceptor.ConnectorConnect(ctx, connect)
}

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
	ctx, span := i.tracer.StartClient(ctx, operationName("exec statement"), attributes.String("db.sql.query", query))
	defer span.End()

	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return res, err
}

func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	ctx, span := i.tracer.StartClient(ctx, operationName("query statement"), attributes.String("db.sql.query", query))
	defer span.End()

	ctx, res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

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

func operationName(name string) string {
	return tracer.OperationName("sql", name)
}
