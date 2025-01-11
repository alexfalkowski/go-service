package retry

import (
	"context"

	"github.com/alexfalkowski/go-service/time"
	"github.com/sethvargo/go-retry"
)

// Try the function.
func Try(ctx context.Context, fn retry.RetryFunc, cfg *Config) error {
	back := retry.NewConstant(time.MustParseDuration(cfg.Backoff))
	back = retry.WithMaxRetries(cfg.Attempts, back)

	return retry.Do(ctx, back, fn)
}
