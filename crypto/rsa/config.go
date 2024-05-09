package rsa

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for rsa.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for rsa.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// GetPrivate from config or env.
func (c Config) GetPrivate() string {
	return os.GetFromEnv(c.Private)
}
