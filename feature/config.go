package feature

import (
	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for feature.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for feature.
type Config struct {
	*client.Config `yaml:",inline" json:",inline" toml:",inline"`
}
