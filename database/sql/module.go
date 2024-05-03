package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	pg.Module,
	fx.Invoke(metrics.Register),
)
