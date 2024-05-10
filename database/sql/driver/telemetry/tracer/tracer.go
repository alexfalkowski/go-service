package tracer

import (
	"context"
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewInterceptor for tracer.
func NewInterceptor(driver string, tracer trace.Tracer, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{driver: driver, tracer: tracer, interceptor: interceptor}
}

// Interceptor for tracer.
type Interceptor struct {
	tracer      trace.Tracer
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
	attrs := []attribute.KeyValue{
		attribute.Key("db.sql.query").String(query),
	}

	ctx, span := i.tracer.Start(ctx, operationName("exec conn"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToString(span.SpanContext().TraceID()))
	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return res, err
}

func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	attrs := []attribute.KeyValue{
		attribute.Key("db.sql.query").String(query),
	}

	ctx, span := i.tracer.Start(ctx, operationName("query conn"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToString(span.SpanContext().TraceID()))
	ctx, res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

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
	attrs := []attribute.KeyValue{
		attribute.Key("db.sql.query").String(query),
	}

	ctx, span := i.tracer.Start(ctx, operationName("exec statement"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToString(span.SpanContext().TraceID()))
	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)

	tracer.Error(err, span)
	tracer.Meta(ctx, span)

	return res, err
}

func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	attrs := []attribute.KeyValue{
		attribute.Key("db.sql.query").String(query),
	}

	ctx, span := i.tracer.Start(ctx, operationName("query statement"), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
	defer span.End()

	ctx = tm.WithTraceID(ctx, meta.ToString(span.SpanContext().TraceID()))
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
