package limiter

// IsEnabled limiter.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for limiter.
type Config struct {
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty" toml:"pattern,omitempty"`
}
