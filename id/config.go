package id

// Config for id.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled for id.
func (c *Config) IsEnabled() bool {
	return c != nil
}
