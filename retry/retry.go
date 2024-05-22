package retry

import (
	"github.com/alexfalkowski/go-service/time"
	"github.com/avast/retry-go/v3"
)

// Try the function.
func Try(fn retry.RetryableFunc, cfg *Config) error {
	d := time.MustParseDuration(cfg.Backoff)

	return retry.Do(fn, retry.Attempts(cfg.Attempts), retry.Delay(d))
}
