package sql

import "github.com/alexfalkowski/go-service/v2/database/sql/pg"

// Config contains SQL database configuration for go-service.
//
// It currently supports PostgreSQL via the PG field.
type Config struct {
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}

// IsEnabled reports whether SQL configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
