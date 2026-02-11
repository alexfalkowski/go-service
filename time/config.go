package time

// Config configures a network time provider.
//
// A nil *Config is treated as disabled (see IsEnabled).
type Config struct {
	// Kind selects the network time provider implementation (for example "ntp" or "nts").
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Address is the provider address passed to the selected implementation.
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

// IsEnabled reports whether the time network provider is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}
