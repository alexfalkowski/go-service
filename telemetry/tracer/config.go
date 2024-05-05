package tracer

// IsEnabled for tracer.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
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
