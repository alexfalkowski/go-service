package config

import (
	"os"
	"path/filepath"
	"strings"
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

	// URL for DSN.
	URL string

	// DSN for SQL.
	DSN struct {
		URL URL `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	}
)

// GetPassword for SQL.
func (d *DSN) GetURL() (string, error) {
	k, err := os.ReadFile(filepath.Clean(string(d.URL)))

	return strings.TrimSpace(string(k)), err
}
