package token

import (
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// Config for token.
type Config struct {
	Access *access.Config `yaml:"access,omitempty" json:"access,omitempty" toml:"access,omitempty"`
	JWT    *jwt.Config    `yaml:"jwt,omitempty" json:"jwt,omitempty" toml:"jwt,omitempty"`
	Paseto *paseto.Config `yaml:"paseto,omitempty" json:"paseto,omitempty" toml:"paseto,omitempty"`
	SSH    *ssh.Config    `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
	Kind   string         `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled for configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
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
