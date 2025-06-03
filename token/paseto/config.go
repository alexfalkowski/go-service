package paseto

// IsEnabled for paseto.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for paseto.
type Config struct {
	Secret     string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
	Issuer     string `yaml:"iss,omitempty" json:"iss,omitempty" toml:"iss,omitempty"`
	Expiration string `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty"`
}
