package paseto

import "github.com/alexfalkowski/go-service/v2/time"

// Config configures PASETO issuance and verification for the go-service PASETO token kind.
//
// This configuration is intended to capture the common knobs for token issuance:
//
//   - Secret: key material configuration (often sourced from env/files)
//   - Issuer: the expected issuer ("iss") value
//   - Expiration: how long issued tokens should be valid
//
// # Secret field note
//
// Although this type includes Secret, the current PASETO implementation in this repository
// (see Token in paseto.go) issues PASETO v4 public tokens using Ed25519 key material supplied
// via crypto/ed25519.Signer and crypto/ed25519.Verifier passed to NewToken. As a result, Secret
// is not consumed directly by that implementation.
//
// If you want to source Ed25519 key material from configuration, resolve Secret using the
// go-service “source string” convention (for example via os.FS.ReadSource) and build the
// Ed25519 signer/verifier in your wiring layer.
//
// # Expiration parsing and panics
//
// Token issuance uses Expiration directly. In config files it is encoded using the standard
// Go duration string format, so invalid values fail during decoding. Apply additional
// validation earlier if you need stricter startup policy.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *Config disables the PASETO implementation and
// NewToken returns nil.
type Config struct {
	// Secret is a "source string" intended to provide PASETO key material.
	//
	// It supports the go-service “source string” pattern:
	// - "env:NAME" to read from an environment variable
	// - "file:/path" to read from a file
	// - otherwise treated as the literal value
	//
	// Note: the current PASETO token implementation in this repository does not read
	// Secret directly; it uses Ed25519 key material provided to NewToken.
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`

	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`

	// Expiration is the duration used to set token expiration.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "24h".
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}

// IsEnabled reports whether PASETO configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
