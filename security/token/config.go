package token

// IsEnabled the config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// Key for token.
	Key string

	// Config for token.
	Config struct {
		Key  Key    `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
		Hash string `yaml:"hash,omitempty" json:"hash,omitempty" toml:"hash,omitempty"`
	}
)
