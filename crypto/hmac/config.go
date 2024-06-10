package hmac

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for hmac.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for hmac.
type Config struct {
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// GetKey for hmac.
func (c *Config) GetKey() (string, error) {
	return os.ReadBase64File(c.Key)
}
