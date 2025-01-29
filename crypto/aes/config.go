package aes

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for aes.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for aes.
type Config struct {
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// GetKey for aes.
func (c *Config) GetKey() (string, error) {
	return os.ReadBase64File(c.Key)
}
