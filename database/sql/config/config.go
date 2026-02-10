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
	ConnMaxLifetime string `yaml:"conn_max_lifetime,omitempty" json:"conn_max_lifetime,omitempty" toml:"conn_max_lifetime,omitempty"`
	Masters         []DSN  `yaml:"masters,omitempty" json:"masters,omitempty" toml:"masters,omitempty"`
	Slaves          []DSN  `yaml:"slaves,omitempty" json:"slaves,omitempty" toml:"slaves,omitempty"`
	MaxOpenConns    int    `yaml:"max_open_conns,omitempty" json:"max_open_conns,omitempty" toml:"max_open_conns,omitempty"`
	MaxIdleConns    int    `yaml:"max_idle_conns,omitempty" json:"max_idle_conns,omitempty" toml:"max_idle_conns,omitempty"`
}

// IsEnabled reports whether SQL configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}

// DSN is a SQL datasource name (connection string) configuration.
type DSN struct {
	// URL is a "source string" for the DSN (literal value, `file:`, or `env:`).
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetURL reads the configured DSN URL bytes from its source.
func (d *DSN) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
