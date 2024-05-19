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

// GetKey for aes.
func (c *Config) GetKey() ([]byte, error) {
	return os.ReadBase64File(string(c.Key))
}
