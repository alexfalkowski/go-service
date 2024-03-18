package tracer

// IsEnabled for tracer.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// Config for tracer.
type Config struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host    string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
	Key     string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Secure  bool   `yaml:"secure,omitempty" json:"secure,omitempty" toml:"secure,omitempty"`
}

// IsBaselime configuration.
func (c *Config) IsBaselime() bool {
	return c.Kind == "baselime"
}
