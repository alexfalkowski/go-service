package retry

// Config configures retry behavior.
type Config struct {
	// Timeout is the per-attempt timeout duration, encoded as a Go duration string (for example "250ms", "5s").
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`

	// Backoff is the backoff duration between attempts, encoded as a Go duration string (for example "100ms", "1s").
	Backoff string `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty"`

	// Attempts is the maximum number of attempts (including the initial attempt).
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}
