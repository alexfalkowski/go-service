package sql

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	pg.Module,
	telemetry.Module,
)
