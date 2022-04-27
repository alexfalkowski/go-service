package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

var (
	// PostgreSQLModule for fx.
	// nolint:gochecknoglobals
	PostgreSQLModule = fx.Provide(pg.NewDB)

	// PostgreSQLJaegerModule for fx.
	// nolint:gochecknoglobals
	PostgreSQLJaegerModule = fx.Provide(jaeger.NewTracer)

	// PostgreSQLDataDogModule for fx.
	// nolint:gochecknoglobals
	PostgreSQLDataDogModule = fx.Provide(datadog.NewTracer)
)
