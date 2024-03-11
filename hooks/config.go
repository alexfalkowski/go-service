package hooks

// Config for hooks.
type Config struct {
	Secret string `yaml:"secret" json:"secret" toml:"secret"`
}
