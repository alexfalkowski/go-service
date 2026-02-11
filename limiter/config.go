package limiter

// Config configures the limiter.
//
// A nil *Config is treated as disabled (see IsEnabled).
type Config struct {
	// Kind selects which limiter key kind is used for limiting (see NewKeyMap for default kinds).
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Interval is a duration string (for example "1s" or "1m") that defines the refill window.
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty"`

	// Tokens is the maximum number of tokens available per interval.
	Tokens uint64 `yaml:"tokens,omitempty" json:"tokens,omitempty" toml:"tokens,omitempty"`
}

// IsEnabled reports whether the limiter is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}
