package retry

// Config configures retry behavior for an operation.
//
// This package defines configuration only; concrete retry behavior is implemented
// by transport-specific packages (for example transport/http/retry and
// transport/grpc/retry). As a result, the exact retry policy (what is retryable,
// whether jitter is applied, etc.) is transport-defined, but these fields are
// the common knobs most implementations interpret similarly.
//
// Timeout and Backoff are encoded as Go duration strings (see time.ParseDuration),
// such as "250ms", "5s", or "1m".
type Config struct {
	// Timeout is the per-attempt timeout duration.
	//
	// When interpreted by a transport, each attempt is typically bounded
	// independently (for example by deriving a per-attempt context deadline).
	//
	// Value encoding: Go duration string (for example "250ms", "5s").
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`

	// Backoff is the delay between attempts after a failed attempt.
	//
	// Transports commonly treat this as a base backoff duration and may apply
	// additional behavior on top (for example jitter), depending on the
	// implementation.
	//
	// Value encoding: Go duration string (for example "100ms", "1s").
	Backoff string `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty"`

	// Attempts is the maximum number of attempts, including the initial attempt.
	//
	// A common convention is:
	//   - Attempts == 1 disables retries (single attempt only).
	//   - Attempts > 1 allows up to Attempts-1 retries after the initial attempt.
	//
	// Some transports may treat Attempts == 0 as "unspecified" and apply a
	// transport default.
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}
