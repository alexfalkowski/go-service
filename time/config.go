package time

// Config for time.
type Config struct {
	Kind    string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

// IsEnabled for time.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// IsNTP for configuration.
func (c *Config) IsNTP() bool {
	return c.Kind == "ntp"
}

// IsNTS for configuration.
func (c *Config) IsNTS() bool {
	return c.Kind == "nts"
}
