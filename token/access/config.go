package access

// Config for access.
type Config struct {
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}

// IsEnabled for access.
func (c *Config) IsEnabled() bool {
	return c != nil
}
