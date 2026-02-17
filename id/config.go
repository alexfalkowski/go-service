package id

// Config configures ID generation.
type Config struct {
	// Kind selects the ID generator implementation (for example "uuid", "ksuid", etc.),
	// depending on which generators are compiled/registered by the service.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled reports whether ID configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
