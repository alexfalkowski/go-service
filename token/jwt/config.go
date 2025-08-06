package jwt

// Config for jwt.
type Config struct {
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
	KeyID      string `yaml:"kid,omitempty" json:"kid,omitempty" toml:"kid,omitempty"`
}

// IsEnabled for jwt.
func (c *Config) IsEnabled() bool {
	return c != nil
}
