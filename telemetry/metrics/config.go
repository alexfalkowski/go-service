package metrics

// IsEnabled for tracer.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for tracer.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
}

// IsOTLP configuration.
func (c *Config) IsOTLP() bool {
	return c.Kind == "otlp"
}
