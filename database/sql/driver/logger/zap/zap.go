package zap

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/ngrok/sqlmw"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	duration  = "%s.duration"
	startTime = "%s.start_time"
	deadline  = "%s.deadline"
	component = "component"
)

// NewInterceptor for zap.
func NewInterceptor(name string, logger *zap.Logger, interceptor sqlmw.Interceptor) *Interceptor {
	return &Interceptor{name: name, logger: logger, interceptor: interceptor}
}

// Interceptor for zap.
type Interceptor struct {
	name        string
	logger      *zap.Logger
	interceptor sqlmw.Interceptor
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

//nolint:dupl
func (i *Interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now().UTC()
	fields := []zapcore.Field{
		zap.String(fmt.Sprintf(startTime, i.name), start.Format(time.RFC3339)),
		zap.String("span.kind", i.name),
		zap.String(component, i.name),
		zap.String(fmt.Sprintf("%s.query", i.name), query),
	}

	for _, a := range args {
		fields = append(fields, zap.Any(fmt.Sprintf("%s.args.%s", i.name, strings.ToLower(a.Name)), a.Value))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(fmt.Sprintf(deadline, i.name), d.UTC().Format(time.RFC3339)))
	}

	res, err := i.interceptor.ConnExecContext(ctx, conn, query, args)

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	fields = append(fields, zap.Int64(fmt.Sprintf(duration, i.name), stime.ToMilliseconds(time.Since(start))))

	if err != nil {
		fields = append(fields, zap.Error(err))
		i.logger.Error("finished call with error", fields...)

		return nil, err
	}

	i.logger.Info("finished call with success", fields...)

	return res, nil
}

func (i *Interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now().UTC()
	fields := []zapcore.Field{
		zap.String(fmt.Sprintf(startTime, i.name), start.Format(time.RFC3339)),
		zap.String("span.kind", i.name),
		zap.String(component, i.name),
		zap.String(fmt.Sprintf("%s.query", i.name), query),
	}

	for _, a := range args {
		fields = append(fields, zap.Any(fmt.Sprintf("%s.args.%s", i.name, strings.ToLower(a.Name)), a.Value))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(fmt.Sprintf(deadline, i.name), d.UTC().Format(time.RFC3339)))
	}

	ctx, res, err := i.interceptor.ConnQueryContext(ctx, conn, query, args)

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	fields = append(fields, zap.Int64(fmt.Sprintf(duration, i.name), stime.ToMilliseconds(time.Since(start))))

	if err != nil {
		fields = append(fields, zap.Error(err))
		i.logger.Error("finished call with error", fields...)

		return ctx, nil, err
	}

	i.logger.Info("finished call with success", fields...)

	return ctx, res, nil
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

//nolint:dupl
func (i *Interceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now().UTC()
	fields := []zapcore.Field{
		zap.String(fmt.Sprintf(startTime, i.name), start.Format(time.RFC3339)),
		zap.String("span.kind", i.name),
		zap.String(component, i.name),
		zap.String(fmt.Sprintf("%s.query", i.name), query),
	}

	for _, a := range args {
		fields = append(fields, zap.Any(fmt.Sprintf("%s.args.%s", i.name, strings.ToLower(a.Name)), a.Value))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(fmt.Sprintf(deadline, i.name), d.UTC().Format(time.RFC3339)))
	}

	res, err := i.interceptor.StmtExecContext(ctx, stmt, query, args)

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	fields = append(fields, zap.Int64(fmt.Sprintf(duration, i.name), stime.ToMilliseconds(time.Since(start))))

	if err != nil {
		fields = append(fields, zap.Error(err))
		i.logger.Error("finished call with error", fields...)

		return nil, err
	}

	i.logger.Info("finished call with success", fields...)

	return res, nil
}

func (i *Interceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	start := time.Now().UTC()
	fields := []zapcore.Field{
		zap.String(fmt.Sprintf(startTime, i.name), start.Format(time.RFC3339)),
		zap.String("span.kind", i.name),
		zap.String(component, i.name),
		zap.String(fmt.Sprintf("%s.query", i.name), query),
	}

	for _, a := range args {
		fields = append(fields, zap.Any(fmt.Sprintf("%s.args.%s", i.name, strings.ToLower(a.Name)), a.Value))
	}

	if d, ok := ctx.Deadline(); ok {
		fields = append(fields, zap.String(fmt.Sprintf(deadline, i.name), d.UTC().Format(time.RFC3339)))
	}

	ctx, res, err := i.interceptor.StmtQueryContext(ctx, stmt, query, args)

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	fields = append(fields, zap.Int64(fmt.Sprintf(duration, i.name), stime.ToMilliseconds(time.Since(start))))

	if err != nil {
		fields = append(fields, zap.Error(err))
		i.logger.Error("finished call with error", fields...)

		return ctx, nil, err
	}

	i.logger.Info("finished call with success", fields...)

	return ctx, res, nil
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
