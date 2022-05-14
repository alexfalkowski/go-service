package opentracing

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/ngrok/sqlmw"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/fx"
)

const (
	pgDuration        = "pg.duration"
	pgStartTime       = "pg.start_time"
	pgRequestDeadline = "pg.request.deadline"
	component         = "component"
	pgComponent       = "pg"
)

// TracerParams for opentracing.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for opentracing.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "pg", Config: params.Config, Version: params.Version})
}

// Tracer for opentracing.
type Tracer otr.Tracer

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "pg", operation, method, opts...)
}

// NewInterceptor for opentracing.
func NewInterceptor(tracer Tracer, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{tracer: tracer, interceptor: interceptor}
}

// Interceptor for opentracing.
type Interceptor struct {
	tracer      Tracer
	interceptor sqlmw.Interceptor
}

func (i *Interceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	return i.interceptor.ConnBeginTx(ctx, conn, txOpts)
}

func (i *Interceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	return i.interceptor.ConnPrepareContext(ctx, conn, query)
}

func (i *Interceptor) ConnPing(ctx context.Context, conn driver.Pinger) error {
	return i.interceptor.ConnPing(ctx, conn)
}

// nolint:dupl
func (i *Interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: pgStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: pgComponent},
		otr.Tag{Key: "pg.query", Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("pg.args.%s", strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := StartSpanFromContext(ctx, i.tracer, "connection", "exec", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(pgRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)
	if err != nil {
		setError(span, err)

		return nil, err
	}

	span.SetTag(pgDuration, stime.ToMilliseconds(time.Since(start)))

	return res, nil
}

// nolint:dupl
func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: pgStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: pgComponent},
		otr.Tag{Key: "pg.query", Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("pg.args.%s", strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := StartSpanFromContext(ctx, i.tracer, "connection", "query", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(pgRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)
	if err != nil {
		setError(span, err)

		return nil, err
	}

	span.SetTag(pgDuration, stime.ToMilliseconds(time.Since(start)))

	return res, nil
}

func (i *Interceptor) ConnectorConnect(ctx context.Context, connect driver.Connector) (driver.Conn, error) {
	return i.interceptor.ConnectorConnect(ctx, connect)
}

// nolint:revive,stylecheck
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

// nolint:dupl
func (i *Interceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: pgStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: pgComponent},
		otr.Tag{Key: "pg.query", Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("pg.args.%s", strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := StartSpanFromContext(ctx, i.tracer, "statement", "exec", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(pgRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)
	if err != nil {
		setError(span, err)

		return nil, err
	}

	span.SetTag(pgDuration, stime.ToMilliseconds(time.Since(start)))

	return res, nil
}

// nolint:dupl
func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now().UTC()
	opts := []otr.StartSpanOption{
		otr.Tag{Key: pgStartTime, Value: start.Format(time.RFC3339)},
		otr.Tag{Key: component, Value: pgComponent},
		otr.Tag{Key: "pg.query", Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("pg.args.%s", strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := StartSpanFromContext(ctx, i.tracer, "statement", "query", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(pgRequestDeadline, d.UTC().Format(time.RFC3339))
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)
	if err != nil {
		setError(span, err)

		return nil, err
	}

	span.SetTag(pgDuration, stime.ToMilliseconds(time.Since(start)))

	return res, nil
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

func setError(span otr.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
}
