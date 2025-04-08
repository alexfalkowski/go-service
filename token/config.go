package token

import (
	"github.com/alexfalkowski/go-service/token/jwt"
	"github.com/alexfalkowski/go-service/token/paseto"
	"github.com/alexfalkowski/go-service/token/ssh"
)

// IsEnabled for token.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for token.
type Config struct {
	JWT    *jwt.Config    `yaml:"jwt,omitempty" json:"jwt,omitempty" toml:"jwt,omitempty"`
	Paseto *paseto.Config `yaml:"paseto,omitempty" json:"paseto,omitempty" toml:"paseto,omitempty"`
	SSH    *ssh.Config    `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
	Kind   string         `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsJWT for configuration.
func (c *Config) IsJWT() bool {
	return c.Kind == "jwt"
}

// IsPaseto for configuration.
func (c *Config) IsPaseto() bool {
	return c.Kind == "paseto"
}

// IsSSH for configuration.
func (c *Config) IsSSH() bool {
	return c.Kind == "ssh"
}
