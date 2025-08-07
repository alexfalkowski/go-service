package sql

import "github.com/alexfalkowski/go-service/v2/database/sql/pg"

// Config for SQL.
type Config struct {
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}

// IsEnabled for SQL.
func (c *Config) IsEnabled() bool {
	return c != nil
}
