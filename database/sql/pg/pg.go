package pg

import (
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/driver"
	"github.com/alexfalkowski/go-service/database/sql/pg/telemetry/tracer"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// OpenParams for pg.
type OpenParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
}

// Open for pg.
//
//nolint:nilnil
func Open(params OpenParams) (*mssqlx.DBs, error) {
	c := params.Config
	if c == nil {
		return nil, nil
	}

	return driver.Open(params.Lifecycle, "pg", c.Config)
}

var once sync.Once

// Register the driver for pg.
func Register(tracer tracer.Tracer, logger *zap.Logger) {
	once.Do(func() {
		driver.Register("pg", stdlib.GetDefaultDriver(), tracer, logger)
	})
}
