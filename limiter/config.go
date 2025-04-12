package limiter

import "github.com/alexfalkowski/go-service/strings"

// IsEnabled limiter.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && !strings.IsEmpty(cfg.Kind)
}

// Config for limiter.
type Config struct {
	Kind     string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty"`
	Tokens   uint64 `yaml:"tokens,omitempty" json:"tokens,omitempty" toml:"tokens,omitempty"`
}
