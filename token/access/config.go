package access

// Config configures access control policy loading.
type Config struct {
	// Policy is a "source string" for the access control policy document.
	//
	// It supports the go-service "source string" pattern:
	// - "env:NAME" to read from an environment variable
	// - "file:/path" to read from a file
	// - otherwise treated as the literal policy content
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}

// IsEnabled reports whether access configuration is present (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}
