package argon2

// IsEnabled the config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for argon2.
type Config struct {
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Hash string `yaml:"hash,omitempty" json:"hash,omitempty" toml:"hash,omitempty"`
}
