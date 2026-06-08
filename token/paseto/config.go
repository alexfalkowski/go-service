package paseto

import (
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/keys"
)

// Config configures PASETO issuance and verification for the go-service PASETO token kind.
//
// This configuration is intended to capture the common knobs for token issuance:
//
//   - Key: active signing key id
//   - Keys: named key material configuration (often sourced from env/files)
//   - Issuer: the expected issuer ("iss") value
//   - Expiration: how long issued tokens should be valid
//
// # Expiration parsing and panics
//
// Token issuance uses Expiration directly. In config files it is encoded using the standard
// Go duration string format, so invalid values fail during decoding. Apply additional
// validation earlier if you need stricter startup policy.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *[Config] disables the PASETO implementation and
// NewToken returns nil.
type Config struct {
	// Key is the active signing key id written to the PASETO footer.
	//
	// The corresponding entry in Keys must include private key material for generation.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty" validate:"required"`

	// Keys contains all named Ed25519 keys trusted for PASETO verification.
	//
	// The key selected by Key is used for signing. Verification reads the token's footer
	// and selects the matching entry from Keys.
	Keys keys.Map `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty" validate:"required"`

	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty" validate:"required"`

	// Expiration is the duration used to set token expiration.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "24h".
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty" validate:"gt=0"`
}

// IsEnabled reports whether PASETO configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
