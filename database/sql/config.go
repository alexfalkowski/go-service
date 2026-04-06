package sql

import "github.com/alexfalkowski/go-service/v2/database/sql/pg"

// Config is the root SQL database configuration for a go-service based service.
//
// It composes driver-specific configuration for SQL backends wired by this repository. Today the
// built-in backend is PostgreSQL via the `PG` field, but the root shape leaves room for additional
// SQL backends without changing the surrounding service config layout.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *Config is treated as "SQL disabled". Driver-specific sub-configs are also pointers and are treated
// as optional; downstream constructors typically return (nil, nil) when a required sub-config is nil.
type Config struct {
	// PG configures PostgreSQL support (master/slave DSNs and shared pool settings).
	//
	// PostgreSQL is the only SQL backend wired by this repository today.
	PG *pg.Config `yaml:"pg,omitempty" json:"pg,omitempty" toml:"pg,omitempty"`
}

// IsEnabled reports whether SQL configuration is present.
//
// By convention, a nil *Config is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}
