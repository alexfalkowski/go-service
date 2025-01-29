package proxy

import (
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for proxy.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for proxy.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
