package token

// IsEnabled the config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for token.
type Config struct {
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Hash string `yaml:"hash,omitempty" json:"hash,omitempty" toml:"hash,omitempty"`
}
