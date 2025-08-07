package config

// Config for cache.
type Config struct {
	Options    map[string]any `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`
	Kind       string         `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Compressor string         `yaml:"compressor,omitempty" json:"compressor,omitempty" toml:"compressor,omitempty"`
	Encoder    string         `yaml:"encoder,omitempty" json:"encoder,omitempty" toml:"encoder,omitempty"`
}

// IsEnabled for cache.
func (c *Config) IsEnabled() bool {
	return c != nil
}
