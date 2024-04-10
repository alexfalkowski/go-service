package metrics

// IsEnabled for tracer.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// Config for tracer.
type Config struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host    string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}
