package limiter

// IsEnabled limiter.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// Config for limiter.
type Config struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty" toml:"pattern,omitempty"`
}
