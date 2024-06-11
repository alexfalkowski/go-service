package token

import (
	"github.com/alexfalkowski/go-service/security/token/argon2"
)

// IsEnabled the config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Argon2 != nil
}

// Config for token.
type Config struct {
	Argon2 *argon2.Config `yaml:"argon2,omitempty" json:"argon2,omitempty" toml:"argon2,omitempty"`
}
