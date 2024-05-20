package hmac

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for hmac.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// Key for hmac.
	Key string

	// Config for hmac.
	Config struct {
		Key Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	}
)

// GetKey for hmac.
func (c *Config) GetKey() (string, error) {
	return os.ReadBase64File(string(c.Key))
}
