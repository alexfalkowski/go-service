package jwt

// IsEnabled for jwt.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for jwt.
type Config struct {
	Subject    string `yaml:"sub,omitempty" json:"sub,omitempty" toml:"sub,omitempty"`
	Audience   string `yaml:"aud,omitempty" json:"aud,omitempty" toml:"aud,omitempty"`
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
	KeyID      string `yaml:"kid,omitempty" json:"kid,omitempty" toml:"kid,omitempty"`
}
