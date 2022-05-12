package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"go.uber.org/fx"
)

var (
	// PostgreSQLModule for fx.
	PostgreSQLModule = fx.Provide(pg.NewDB)

	// PostgreSQLOpentracingModule for fx.
	PostgreSQLOpentracingModule = fx.Provide(opentracing.NewTracer)
)
