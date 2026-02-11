package paseto

// Config configures PASETO issuance and verification.
type Config struct {
	// Secret is a "source string" used as the PASETO key material.
	//
	// It supports the go-service "source string" pattern:
	// - "env:NAME" to read from an environment variable
	// - "file:/path" to read from a file
	// - otherwise treated as the literal value
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`

	// Issuer is written to and verified against the `iss` claim.
	Issuer string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`

	// Expiration is a duration string used to set/validate the token expiration (for example "15m", "24h").
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}

// IsEnabled reports whether PASETO configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
