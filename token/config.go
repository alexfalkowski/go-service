package token

import (
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// Config configures token generation and verification for a go-service application.
//
// This type is typically embedded into a larger service configuration and consumed
// by the top-level token facade (see [NewToken]), which delegates to the
// configured token implementation.
//
// # Enablement model
//
// Enablement is modeled by presence:
//   - A nil *[Config] means "token support disabled" at the top level.
//   - When enabled, Kind must select a supported implementation and that
//     implementation's nested configuration must be present.
//
// # Selecting an implementation (Kind)
//
// Kind selects the token implementation used by the token facade. Supported kinds are:
//
//   - "jwt": JSON Web Tokens (see package token/jwt)
//   - "paseto": PASETO v4 public tokens (see package token/paseto)
//   - "ssh": SSH-style signed tokens (see package token/ssh)
//
// The selected implementation's nested configuration must be provided in the
// corresponding field (JWT/Paseto/SSH).
//
// Config validation rejects unknown Kind values. If validation is bypassed, the
// token facade treats the configuration as invalid and [Token.Generate]/[Token.Verify]
// return [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
type Config struct {
	// JWT configures the JWT token implementation.
	//
	// When Kind == "jwt", this configuration is consumed by token/jwt.
	JWT *jwt.Config `yaml:"jwt,omitempty" json:"jwt,omitempty" toml:"jwt,omitempty" validate:"required_if=Kind jwt"`

	// Paseto configures the PASETO token implementation.
	//
	// When Kind == "paseto", this configuration is consumed by token/paseto.
	Paseto *paseto.Config `yaml:"paseto,omitempty" json:"paseto,omitempty" toml:"paseto,omitempty" validate:"required_if=Kind paseto"`

	// SSH configures the SSH token implementation.
	//
	// When Kind == "ssh", this configuration is consumed by token/ssh.
	SSH *ssh.Config `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty" validate:"required_if=Kind ssh"`

	// Kind selects the token implementation to use.
	//
	// Supported values: "jwt", "paseto", "ssh".
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty" validate:"required,oneof=jwt paseto ssh"`
}

// IsEnabled reports whether token configuration is present.
//
// A nil receiver is considered disabled and returns false. This is commonly used
// as a simple top-level enable/disable switch for token wiring.
func (c *Config) IsEnabled() bool {
	return c != nil
}
