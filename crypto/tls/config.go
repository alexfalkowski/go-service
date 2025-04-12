package tls

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
)

// IsEnabled for security.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for tls.
type Config struct {
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// HasKeyPair for security.
func (c *Config) HasKeyPair() bool {
	return !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// GetCert for tls.
func (c *Config) GetCert() ([]byte, error) {
	return os.ReadFile(c.Cert)
}

// GetKey for tls.
func (c *Config) GetKey() ([]byte, error) {
	return os.ReadFile(c.Key)
}
