package ed25519

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for ed25519.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for ed25519.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// GetPrivate from config or env.
func (c Config) GetPrivate() string {
	return os.GetFromEnv(c.Private)
}
