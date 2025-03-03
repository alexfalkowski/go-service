package aes

import "github.com/alexfalkowski/go-service/os"

// IsEnabled for aes.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for aes.
type Config struct {
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// GetKey for aes.
func (c *Config) GetKey() ([]byte, error) {
	return os.ReadFile(c.Key)
}
