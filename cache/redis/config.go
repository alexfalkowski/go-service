package redis

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for redis.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// URL for redis.
	URL string

	// Config for redis.
	Config struct {
		Compressor string            `yaml:"compressor,omitempty" json:"compressor,omitempty" toml:"compressor,omitempty"`
		Marshaller string            `yaml:"marshaller,omitempty" json:"marshaller,omitempty" toml:"marshaller,omitempty"`
		Addresses  map[string]string `yaml:"addresses,omitempty" json:"addresses,omitempty" toml:"addresses,omitempty"`
		URL        URL               `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`
	}
)

// GetURL for redis.
func (c *Config) GetURL() (string, error) {
	return os.ReadFile(string(c.URL))
}
