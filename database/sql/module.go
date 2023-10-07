package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	"github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"go.uber.org/fx"
)

// PostgreSQLModule for fx.
var PostgreSQLModule = fx.Options(
	fx.Provide(pg.Open),
	fx.Invoke(pg.Register),
	fx.Invoke(metrics.Register),
	fx.Provide(tracer.NewTracer),
)
