package token

import (
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// Config configures token generation and verification for a go-service application.
//
// This type is typically embedded into a larger service configuration and consumed
// by the top-level token facade (see token.NewToken), which delegates to the
// configured token implementation.
//
// # Enablement model
//
// Enablement is modeled by presence:
//   - A nil *Config means "token support disabled" at the top level.
//   - When enabled, individual token implementations may still be disabled if their
//     nested configuration is nil (for example JWT == nil while Kind == "jwt").
//
// # Selecting an implementation (Kind)
//
// Kind selects the token implementation used by the token facade. Supported kinds are:
//
//   - "jwt": JSON Web Tokens (see package token/jwt)
//   - "paseto": PASETO v4 public tokens (see package token/paseto)
//   - "ssh": SSH-style signed tokens (see package token/ssh)
//
// The selected implementation’s nested configuration should typically be provided
// in the corresponding field (JWT/Paseto/SSH).
//
// If Kind is unknown, the token facade intentionally behaves like "disabled":
// Generate returns (nil, nil) and Verify returns (strings.Empty, nil). Callers
// should treat a nil/empty successful result as a signal to check configuration.
//
// # Access control (Access)
//
// Access configures optional access-control policy wiring (see package token/access).
// It is orthogonal to the token kind: some services may use token verification to
// establish identity (subject) and then evaluate permissions via Access.
type Config struct {
	// Access configures access control policy used by token/access.
	//
	// This is used to answer authorization checks (for example "user has permission X")
	// and is typically used in addition to authentication (token verification).
	Access *access.Config `yaml:"access,omitempty" json:"access,omitempty" toml:"access,omitempty"`

	// JWT configures the JWT token implementation.
	//
	// When Kind == "jwt", this configuration is consumed by token/jwt.
	JWT *jwt.Config `yaml:"jwt,omitempty" json:"jwt,omitempty" toml:"jwt,omitempty"`

	// Paseto configures the PASETO token implementation.
	//
	// When Kind == "paseto", this configuration is consumed by token/paseto.
	Paseto *paseto.Config `yaml:"paseto,omitempty" json:"paseto,omitempty" toml:"paseto,omitempty"`

	// SSH configures the SSH token implementation.
	//
	// When Kind == "ssh", this configuration is consumed by token/ssh.
	SSH *ssh.Config `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`

	// Kind selects the token implementation to use.
	//
	// Supported values: "jwt", "paseto", "ssh".
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled reports whether token configuration is present.
//
// A nil receiver is considered disabled and returns false. This is commonly used
// as a simple top-level enable/disable switch for token wiring.
func (c *Config) IsEnabled() bool {
	return c != nil
}
