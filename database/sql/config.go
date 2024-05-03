package sql

import (
	"github.com/alexfalkowski/go-service/database/sql/pg"
)

// IsEnabled for SQL.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for SQL.
type Config struct {
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}
