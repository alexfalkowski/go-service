package tracer

// IsEnabled for tracer.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for tracer.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsBaselime configuration.
func (c *Config) IsBaselime() bool {
	return c.Kind == "baselime"
}
