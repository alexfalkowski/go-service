package opentracing

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/ngrok/sqlmw"
	otr "github.com/opentracing/opentracing-go"
)

const (
	deadline  = "%s.deadline"
	component = "component"
)

// NewInterceptor for opentracing.
func NewInterceptor(driver string, tracer otr.Tracer, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{driver: driver, tracer: tracer, interceptor: interceptor}
}

// Interceptor for opentracing.
type Interceptor struct {
	driver      string
	tracer      otr.Tracer
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
	opts := []otr.StartSpanOption{
		otr.Tag{Key: component, Value: i.driver},
		otr.Tag{Key: fmt.Sprintf("%s.query", i.driver), Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("%s.args.%s", i.driver, strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := opentracing.StartSpanFromContext(ctx, i.tracer, i.driver, "connection", "exec", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(fmt.Sprintf(deadline, i.driver), d.UTC().Format(time.RFC3339))
	}

	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return res, err
}

// nolint:dupl
func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	opts := []otr.StartSpanOption{
		otr.Tag{Key: component, Value: i.driver},
		otr.Tag{Key: fmt.Sprintf("%s.query", i.driver), Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("%s.args.%s", i.driver, strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := opentracing.StartSpanFromContext(ctx, i.tracer, i.driver, "connection", "query", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(fmt.Sprintf(deadline, i.driver), d.UTC().Format(time.RFC3339))
	}

	res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return res, err
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
	opts := []otr.StartSpanOption{
		otr.Tag{Key: component, Value: i.driver},
		otr.Tag{Key: fmt.Sprintf("%s.query", i.driver), Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("%s.args.%s", i.driver, strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := opentracing.StartSpanFromContext(ctx, i.tracer, i.driver, "statement", "exec", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(fmt.Sprintf(deadline, i.driver), d.UTC().Format(time.RFC3339))
	}

	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return res, err
}

// nolint:dupl
func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	opts := []otr.StartSpanOption{
		otr.Tag{Key: component, Value: i.driver},
		otr.Tag{Key: fmt.Sprintf("%s.query", i.driver), Value: query},
	}

	for _, a := range args {
		opts = append(opts, otr.Tag{Key: fmt.Sprintf("%s.args.%s", i.driver, strings.ToLower(a.Name)), Value: a.Value})
	}

	ctx, span := opentracing.StartSpanFromContext(ctx, i.tracer, i.driver, "statement", "query", opts...)
	defer span.Finish()

	if d, ok := ctx.Deadline(); ok {
		span.SetTag(fmt.Sprintf(deadline, i.driver), d.UTC().Format(time.RFC3339))
	}

	res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return res, err
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
