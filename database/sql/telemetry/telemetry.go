package telemetry

import (
	"database/sql"
	"database/sql/driver"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Option is an alias of otelsql.Option.
//
// It configures SQL OpenTelemetry instrumentation behavior such as emitted
// attributes and metric/tracing options.
type Option = otelsql.Option

// SpanOptions is an alias of otelsql.SpanOptions.
//
// It configures SQL tracing span behavior.
type SpanOptions = otelsql.SpanOptions

// Open opens a `database/sql` DB handle with OpenTelemetry instrumentation.
//
// Raw SQL query text capture is disabled by default. Callers that need raw
// statements in spans may opt in with WithSpanOptions.
func Open(driverName, dataSourceName string, options ...Option) (*sql.DB, error) {
	return otelsql.Open(driverName, dataSourceName, optionsWithDefaults(options)...)
}

// WrapDriver wraps a `database/sql/driver.Driver` with OpenTelemetry
// instrumentation.
//
// Raw SQL query text capture is disabled by default. Callers that need raw
// statements in spans may opt in with WithSpanOptions.
func WrapDriver(driver driver.Driver, opts ...Option) driver.Driver {
	return otelsql.WrapDriver(driver, optionsWithDefaults(opts)...)
}

// WithAttributes adds static attributes to SQL telemetry spans and metrics.
//
// This is a thin wrapper around otelsql.WithAttributes.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return otelsql.WithAttributes(attrs...)
}

// WithSpanOptions configures SQL tracing span behavior.
//
// This is a thin wrapper around otelsql.WithSpanOptions.
func WithSpanOptions(opts SpanOptions) Option {
	return otelsql.WithSpanOptions(opts)
}

// RegisterDBStatsMetrics registers OpenTelemetry DB stats metrics for db.
//
// This is a thin wrapper around otelsql.RegisterDBStatsMetrics.
func RegisterDBStatsMetrics(db *sql.DB, opts ...Option) (metric.Registration, error) {
	return otelsql.RegisterDBStatsMetrics(db, opts...)
}

func optionsWithDefaults(options []Option) []Option {
	opts := make([]Option, 0, len(options)+1)
	opts = append(opts, WithSpanOptions(SpanOptions{DisableQuery: true}))
	opts = append(opts, options...)

	return opts
}
