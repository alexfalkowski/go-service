package ntp

// IsEnabled for ntp.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for ntp.
type Config struct {
	Host string `yaml:"host,omitempty" json:"host,omitempty" toml:"host,omitempty"`
}
