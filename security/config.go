package security

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for security.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// Config for security.
type Config struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Cert    string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`
	Key     string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// HasKeyPair for security.
func (c *Config) HasKeyPair() bool {
	return c.GetCert() != "" && c.GetKey() != ""
}

// GetCert for security.
func (c *Config) GetCert() string {
	return os.GetFromEnv(c.Cert)
}

// GetKey for security.
func (c *Config) GetKey() string {
	return os.GetFromEnv(c.Key)
}
