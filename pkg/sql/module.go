package sql

import (
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	"go.uber.org/fx"
)

// PostgreSQLModule for fx.
// nolint:gochecknoglobals
var PostgreSQLModule = fx.Options(fx.Provide(pg.NewDB))
