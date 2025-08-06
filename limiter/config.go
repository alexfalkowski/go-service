package limiter

// Config for limiter.
type Config struct {
	Kind     string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty"`
	Tokens   uint64 `yaml:"tokens,omitempty" json:"tokens,omitempty" toml:"tokens,omitempty"`
}

// IsEnabled for limiter.
func (c *Config) IsEnabled() bool {
	return c != nil
}
