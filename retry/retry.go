package retry

import (
	"github.com/avast/retry-go/v3"
)

// Try the function.
func Try(fn retry.RetryableFunc, cfg *Config) error {
	return retry.Do(fn, retry.Attempts(cfg.Attempts))
}
