package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for SQL.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for SQL.
type Config struct {
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}
