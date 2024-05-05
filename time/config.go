package time

// IsEnabled for time.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Kind != ""
}

// Config for time.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
}

// IsNTP for configuration.
func (c *Config) IsNTP() bool {
	return c.Kind == "ntp"
}

// IsNTS for configuration.
func (c *Config) IsNTS() bool {
	return c.Kind == "nts"
}
