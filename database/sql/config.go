package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
)

// Config for sql.
type Config struct {
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}
