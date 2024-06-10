package hooks

import (
	"github.com/alexfalkowski/go-service/os"
)

// Config for hooks.
type Config struct {
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}

// GetCert for hooks.
func (c *Config) GetSecret() (string, error) {
	return os.ReadFile(c.Secret)
}
