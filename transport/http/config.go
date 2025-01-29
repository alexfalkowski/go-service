package http

import (
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for HTTP.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for HTTP.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
