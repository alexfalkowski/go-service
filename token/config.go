package token

// IsEnabled for token.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for token.
type Config struct {
	Kind       string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Key        string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Subject    string `yaml:"sub,omitempty" json:"sub,omitempty" toml:"sub,omitempty"`
	Audience   string `yaml:"aud,omitempty" json:"aud,omitempty" toml:"aud,omitempty"`
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}

// IsPaseto for configuration.
func (c *Config) IsKey() bool {
	return c.Kind == "key"
}

// IsJWT for configuration.
func (c *Config) IsJWT() bool {
	return c.Kind == "jwt"
}

// IsPaseto for configuration.
func (c *Config) IsPaseto() bool {
	return c.Kind == "paseto"
}
