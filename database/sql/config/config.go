package config

import (
	"time"
)

// Config for SQL DB.
type Config struct {
	Masters         []DSN         `yaml:"masters" json:"masters" toml:"masters"`
	Slaves          []DSN         `yaml:"slaves" json:"slaves" toml:"slaves"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns" toml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns" toml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime" toml:"conn_max_lifetime"`
}

// DSN for SQL DB.
type DSN struct {
	URL string `yaml:"url" json:"url" toml:"url"`
}
