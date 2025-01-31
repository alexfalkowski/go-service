package cache

// IsEnabled for cache.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for cache.
type Config struct {
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`
	Kind    string         `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}
