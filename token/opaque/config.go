package opaque

// IsEnabled for opaque.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && cfg.Name != ""
}

// Config for opaque.
type Config struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
}
