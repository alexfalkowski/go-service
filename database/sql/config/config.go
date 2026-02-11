package config

import "github.com/alexfalkowski/go-service/v2/os"

// Config contains shared `database/sql` pool configuration.
//
// It is intended to be embedded by driver-specific configuration types (for
// example, PostgreSQL) and consumed by the SQL wiring in this repository.
//
// Masters and Slaves contain DSNs (connection strings) expressed as "source
// strings" (literal values, `file:` paths, or `env:` references).
type Config struct {
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// The value is parsed as a Go duration string (for example "30s", "5m", "1h").
	ConnMaxLifetime string `yaml:"conn_max_lifetime,omitempty" json:"conn_max_lifetime,omitempty" toml:"conn_max_lifetime,omitempty"`

	// Masters is the set of primary (read-write) datasource DSNs.
	//
	// Each DSN URL is a "source string" (literal value, `file:`, or `env:`).
	Masters []DSN `yaml:"masters,omitempty" json:"masters,omitempty" toml:"masters,omitempty"`

	// Slaves is the set of replica (read-only) datasource DSNs.
	//
	// Each DSN URL is a "source string" (literal value, `file:`, or `env:`).
	Slaves []DSN `yaml:"slaves,omitempty" json:"slaves,omitempty" toml:"slaves,omitempty"`

	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `yaml:"max_open_conns,omitempty" json:"max_open_conns,omitempty" toml:"max_open_conns,omitempty"`

	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	MaxIdleConns int `yaml:"max_idle_conns,omitempty" json:"max_idle_conns,omitempty" toml:"max_idle_conns,omitempty"`
}

// IsEnabled reports whether SQL configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}

// DSN is a SQL datasource name (connection string) configuration.
type DSN struct {
	// URL is a "source string" for the datasource name/connection string.
	//
	// It supports the go-service "source string" pattern:
	// - "env:NAME" to read from an environment variable
	// - "file:/path" to read from a file
	// - otherwise treated as a literal connection string
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetURL reads the configured DSN URL bytes from its source.
func (d *DSN) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
