package hooks

import (
	"github.com/alexfalkowski/go-service/os"
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
func (c *Config) GetSecret() (string, error) {
	return os.ReadFile(string(c.Secret))
}
