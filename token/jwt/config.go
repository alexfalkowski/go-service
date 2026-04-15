package jwt

import "github.com/alexfalkowski/go-service/v2/time"

// Config configures JWT issuance and verification for the go-service JWT token kind.
//
// This configuration is consumed by Token.Generate and Token.Verify.
//
// # Claims and headers
//
// Issued tokens use standard registered claims and populate them from this config:
//
//   - Issuer is written to and verified against the `iss` claim.
//   - Expiration is a typed duration and is used to set/validate the `exp` claim.
//   - KeyID is written to and verified against the JWT header `kid`.
//
// # Key ID (kid) enforcement
//
// Verification in this repository is intentionally strict about the `kid` header:
//
//   - The header must exist and be non-empty.
//   - The value must match KeyID exactly.
//
// This is part of the verification contract and helps prevent accepting tokens minted for a
// different key identity.
//
// # Expiration parsing and panics
//
// Issuance uses Expiration directly. Invalid config-file values fail during decoding because
// the field is encoded as a Go duration string. If your deployment requires additional
// validation policy, apply it earlier during configuration loading.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *Config disables the JWT implementation and
// NewToken returns nil.
type Config struct {
	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`

	// KeyID is written to and verified against the JWT header `kid`.
	//
	// Note: this repository's JWT verification expects the `kid` header to be set and
	// to match this value exactly.
	KeyID string `yaml:"kid,omitempty" json:"kid,omitempty" toml:"kid,omitempty"`

	// Expiration is the duration used to set and validate the `exp` claim.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "24h".
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}

// IsEnabled reports whether JWT configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
