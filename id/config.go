package id

import "github.com/alexfalkowski/go-service/strings"

// IsEnabled the config.
func IsEnabled(config *Config) bool {
	return config != nil && !strings.IsEmpty(config.Kind)
}

// Config for id.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}
