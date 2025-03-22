package id

// IsEnabled the config.
func IsEnabled(config *Config) bool {
	return config != nil && config.Kind != ""
}

// Config for id.
type Config struct {
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}
