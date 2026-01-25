package config

import "github.com/alexfalkowski/go-service/v2/os"

// Config for SQL.
type Config struct {
	Name            string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime,omitempty" json:"conn_max_lifetime,omitempty" toml:"conn_max_lifetime,omitempty"`
	Masters         []DSN  `yaml:"masters,omitempty" json:"masters,omitempty" toml:"masters,omitempty"`
	Slaves          []DSN  `yaml:"slaves,omitempty" json:"slaves,omitempty" toml:"slaves,omitempty"`
	MaxOpenConns    int    `yaml:"max_open_conns,omitempty" json:"max_open_conns,omitempty" toml:"max_open_conns,omitempty"`
	MaxIdleConns    int    `yaml:"max_idle_conns,omitempty" json:"max_idle_conns,omitempty" toml:"max_idle_conns,omitempty"`
}

// IsEnabled for SQL.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// DSN for SQL.
type DSN struct {
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
}

// GetURL for SQL.
func (d *DSN) GetURL(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(d.URL)
}
