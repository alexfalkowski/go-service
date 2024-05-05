package limiter

// IsEnabled limiter.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for limiter.
type Config struct {
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty" toml:"pattern,omitempty"`
}
