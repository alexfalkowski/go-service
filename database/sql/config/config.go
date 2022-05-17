package config

import (
	"time"
)

// Config for SQL DB.
type Config struct {
	Masters         []DSN         `yaml:"masters"`
	Slaves          []DSN         `yaml:"slaves"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// DSN for SQL DB.
type DSN struct {
	URL string `yaml:"url"`
}
