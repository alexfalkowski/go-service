package token

import (
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// Config contains token configuration for go-service.
//
// Kind selects which token implementation is used by the top-level token helper
// (for example: "jwt", "paseto", or "ssh").
// The corresponding nested configuration for the chosen kind should also be provided.
type Config struct {
	Access *access.Config `yaml:"access,omitempty" json:"access,omitempty" toml:"access,omitempty"`
	JWT    *jwt.Config    `yaml:"jwt,omitempty" json:"jwt,omitempty" toml:"jwt,omitempty"`
	Paseto *paseto.Config `yaml:"paseto,omitempty" json:"paseto,omitempty" toml:"paseto,omitempty"`
	SSH    *ssh.Config    `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
	Kind   string         `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled reports whether token configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
