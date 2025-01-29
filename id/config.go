package id

import "github.com/alexfalkowski/go-service/structs"

// IsEnabled the config.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for id.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}
