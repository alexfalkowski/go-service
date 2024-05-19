package hooks

type (
	// Secret for hooks.
	Secret string

	// Config for hooks.
	Config struct {
		Secret Secret `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
	}
)
