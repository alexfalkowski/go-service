package paseto

// Config configures PASETO issuance and verification.
//
// Secret is a "source string" (literal value, `file:`, or `env:`) used as the key material.
// Issuer is written to and verified against the `iss` claim.
// Expiration is a duration string used to set the token expiration.
type Config struct {
	Secret     string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}

// IsEnabled reports whether PASETO configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
