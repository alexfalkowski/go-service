package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"go.uber.org/fx"
)

var (
	// PostgreSQLModule for fx.
	PostgreSQLModule = fx.Options(
		fx.Provide(pg.Open),
		fx.Invoke(pg.Register),
	)

	// PostgreSQLOpentracingModule for fx.
	PostgreSQLOpentracingModule = fx.Provide(opentracing.NewTracer)
)
