package retry

// Config for retry.
type Config struct {
	Timeout  string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Attempts uint   `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}
