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

// WrapDriver wraps a `database/sql/driver.Driver` with OpenTelemetry
// instrumentation.
//
// This is a thin wrapper around otelsql.WrapDriver.
func WrapDriver(driver driver.Driver, opts ...Option) driver.Driver {
	return otelsql.WrapDriver(driver, opts...)
}

// WithAttributes adds static attributes to SQL telemetry spans and metrics.
//
// This is a thin wrapper around otelsql.WithAttributes.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return otelsql.WithAttributes(attrs...)
}

// RegisterDBStatsMetrics registers OpenTelemetry DB stats metrics for db.
//
// This is a thin wrapper around otelsql.RegisterDBStatsMetrics.
func RegisterDBStatsMetrics(db *sql.DB, opts ...Option) (metric.Registration, error) {
	return otelsql.RegisterDBStatsMetrics(db, opts...)
}
