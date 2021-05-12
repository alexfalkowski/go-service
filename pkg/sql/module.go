package sql

import (
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	"go.uber.org/fx"
)

var (
	// PostgreSQLModule for fx.
	PostgreSQLModule = fx.Options(fx.Provide(pg.NewDB), fx.Provide(pg.NewConfig))
)
