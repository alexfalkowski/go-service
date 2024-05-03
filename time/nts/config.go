package nts

// IsEnabled for nts.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for nts.
type Config struct {
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
}
