package config

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for SQL.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// Config for SQL.
	Config struct {
		ConnMaxLifetime string `yaml:"conn_max_lifetime,omitempty" json:"conn_max_lifetime,omitempty" toml:"conn_max_lifetime,omitempty"`
		Masters         []DSN  `yaml:"masters,omitempty" json:"masters,omitempty" toml:"masters,omitempty"`
		Slaves          []DSN  `yaml:"slaves,omitempty" json:"slaves,omitempty" toml:"slaves,omitempty"`
		MaxOpenConns    int    `yaml:"max_open_conns,omitempty" json:"max_open_conns,omitempty" toml:"max_open_conns,omitempty"`
		MaxIdleConns    int    `yaml:"max_idle_conns,omitempty" json:"max_idle_conns,omitempty" toml:"max_idle_conns,omitempty"`
	}

	// DSN for SQL.
	DSN struct {
		URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	}
)

// GetPassword for SQL.
func (d *DSN) GetURL() (string, error) {
	return os.ReadFile(d.URL)
}
