package hooks

// Config for hooks.
type Config struct {
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}
