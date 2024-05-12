package rsa

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for rsa.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// PublicKey for rsa.
	PublicKey string

	// PrivateKey for rsa.
	PrivateKey string

	// Config for rsa.
	Config struct {
		Public  PublicKey  `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
		Private PrivateKey `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
	}
)

// GetPrivate from config or env.
func (c *Config) GetPrivate() string {
	return os.GetFromEnv(string(c.Private))
}
