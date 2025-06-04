package access

// IsEnabled for access.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for access.
type Config struct {
	Model  string `yaml:"model,omitempty" json:"model,omitempty" toml:"model,omitempty"`
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}
