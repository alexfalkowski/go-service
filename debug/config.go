package debug

import (
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for debug.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for debug.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
