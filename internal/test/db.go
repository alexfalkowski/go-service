package test

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/linxGnu/mssqlx"
)

// WithWorldPGConfig overrides the Postgres config used when the world opens databases.
func WithWorldPGConfig(config *pg.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.pg = config
	})
}

// OpenDatabase opens the world's configured Postgres connection set through the shared test filesystem.
func (w *World) OpenDatabase() (*mssqlx.DBs, error) {
	dbs, err := pg.Open(w.Lifecycle, FS, w.PG)
	if err != nil {
		return nil, err
	}

	return dbs, err
}

func (w *World) registerDatabase() {
	pg.Register()
}

func pgConfig(os *worldOpts) *pg.Config {
	if os.pg != nil {
		return os.pg
	}

	return NewPGConfig()
}
