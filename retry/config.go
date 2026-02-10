package retry

// Config configures retry behavior.
//
// Timeout is the per-attempt timeout duration.
// Backoff is the backoff duration between attempts.
// Attempts is the maximum number of attempts (including the initial attempt).
type Config struct {
	Timeout  string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
	Backoff  string `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty"`
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}
