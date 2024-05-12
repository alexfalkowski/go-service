package aes

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for aes.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// Key for aes.
	Key string

	// Config for aes.
	Config struct {
		Key Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	}
)

// GetKey from config or env.
func (c *Config) GetKey() string {
	return os.GetFromEnv(string(c.Key))
}
