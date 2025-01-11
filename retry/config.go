package retry

// Config for retry.
type Config struct {
	Timeout  string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Backoff  string `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty"`
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}
