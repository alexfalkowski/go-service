package jwt

import (
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/keys"
)

// Config configures JWT issuance and verification for the go-service JWT token kind.
//
// This configuration is consumed by [Token.Generate] and [Token.Verify].
//
// # Claims and headers
//
// Issued tokens use standard registered claims and populate them from this config:
//
//   - Issuer is written to and verified against the `iss` claim.
//   - Expiration is a typed duration and is used to set/validate the `exp` claim.
//   - Key is written to the JWT header `kid`.
//
// # Key ID (kid) enforcement
//
// Verification in this repository is intentionally strict about the `kid` header:
//
//   - The header must exist and be non-empty.
//   - The value must select a configured entry from Keys.
//
// This is part of the verification contract and helps prevent accepting tokens minted for a
// different key identity.
//
// # Validation
//
// Issuer, Key, and Keys are required because generated tokens must be verifiable by
// this package. Expiration must be greater than zero because zero-duration tokens
// are immediately expired.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *[Config] disables the JWT implementation and
// NewToken returns nil.
type Config struct {
	// Keys contains all named Ed25519 keys trusted for JWT verification.
	//
	// The key selected by Key is used for signing. Verification reads the token's `kid`
	// header and selects the matching entry from Keys.
	Keys keys.Map `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty" validate:"required"`

	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty" validate:"required"`

	// Key is the active signing key id written to the JWT header `kid`.
	//
	// The corresponding entry in Keys must include private key material for generation.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty" validate:"required"`

	// Expiration is the duration used to set and validate the `exp` claim.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "24h".
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty" validate:"gt=0"`
}

// IsEnabled reports whether JWT configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
