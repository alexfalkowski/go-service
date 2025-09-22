package test

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry/metrics"
	"github.com/linxGnu/mssqlx"
)

// WithWorldPGConfig for test.
func WithWorldPGConfig(config *pg.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.pg = config
	})
}

// OpenDatabase for world.
func (w *World) OpenDatabase() (*mssqlx.DBs, error) {
	dbs, err := pg.Open(w.Lifecycle, FS, w.PG)
	if err != nil {
		return nil, err
	}

	metrics.Register(dbs, w.Server.Meter)

	return dbs, err
}

func (w *World) registerDatabase() {
	pg.Register(w.NewTracer(), w.Logger)
}

func pgConfig(os *worldOpts) *pg.Config {
	if os.pg != nil {
		return os.pg
	}

	return NewPGConfig()
}
