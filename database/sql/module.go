package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"go.uber.org/fx"
)

// PostgreSQLModule for fx.
// nolint:gochecknoglobals
var PostgreSQLModule = fx.Options(fx.Provide(pg.NewDB))
