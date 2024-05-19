package tls

import (
	"os"
	"path/filepath"
)

// IsEnabled for security.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (

	// Cert for tls.
	Cert string

	// Cert for tls.
	Key string

	// Config for tls.
	Config struct {
		Cert Cert `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`
		Key  Key  `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	}
)

// HasKeyPair for security.
func (c *Config) HasKeyPair() bool {
	return c.Cert != "" && c.Key != ""
}

// GetCert for tls.
func (c *Config) GetCert() ([]byte, error) {
	return os.ReadFile(filepath.Clean(string(c.Cert)))
}

// GetKey for tls.
func (c *Config) GetKey() ([]byte, error) {
	return os.ReadFile(filepath.Clean(string(c.Key)))
}
