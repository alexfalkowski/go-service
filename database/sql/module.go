package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/otel"
	"go.uber.org/fx"
)

var (
	// PostgreSQLModule for fx.
	PostgreSQLModule = fx.Options(
		fx.Provide(pg.Open),
		fx.Invoke(pg.Register),
	)

	// PostgreSQLOTELModule for fx.
	PostgreSQLOTELModule = fx.Provide(otel.NewTracer)
)
