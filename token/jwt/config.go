package jwt

// Config configures JWT issuance and verification.
//
// Issuer is written to and verified against the `iss` claim.
// Expiration is a duration string used to set the `exp` claim.
// KeyID is written to and verified against the JWT header `kid`.
type Config struct {
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
	KeyID      string `yaml:"kid,omitempty" json:"kid,omitempty" toml:"kid,omitempty"`
}

// IsEnabled reports whether JWT configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
