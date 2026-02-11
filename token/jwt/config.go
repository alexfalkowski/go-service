package jwt

// Config configures JWT issuance and verification.
type Config struct {
	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`

	// Expiration is a duration string used to set/validate the `exp` claim (for example "15m", "24h").
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`

	// KeyID is written to and verified against the JWT header `kid`.
	//
	// Note: this repository’s JWT verification expects the `kid` header to be set.
	KeyID string `yaml:"kid,omitempty" json:"kid,omitempty" toml:"kid,omitempty"`
}

// IsEnabled reports whether JWT configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
