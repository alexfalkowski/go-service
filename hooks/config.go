package hooks

import (
	"os"
	"path/filepath"
)

type (
	// Secret for hooks.
	Secret string

	// Config for hooks.
	Config struct {
		Secret Secret `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
	}
)

// GetCert for hooks.
func (c *Config) GetSecret() ([]byte, error) {
	return os.ReadFile(filepath.Clean(string(c.Secret)))
}
