package retry

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	config "github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	retryable "github.com/hashicorp/go-retryablehttp"
	"github.com/sethvargo/go-retry"
)

// Config is an alias for retry.Config.
type Config = config.Config

// ErrInvalidStatusCode for http retry.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

// NewRoundTripper for retry.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper) *RoundTripper {
	timeout := time.MustParseDuration(cfg.Timeout)
	backoff := retry.WithMaxRetries(cfg.Attempts, retry.NewConstant(time.MustParseDuration(cfg.Backoff)))

	return &RoundTripper{RoundTripper: hrt, timeout: timeout, backoff: backoff}
}

// RoundTripper for retry.
type RoundTripper struct {
	http.RoundTripper
	backoff retry.Backoff
	timeout time.Duration
}

// RoundTrip executes the request with retries according to the config.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	operation := func(ctx context.Context) (*http.Response, error) {
		ctx, cancel := context.WithTimeout(ctx, r.timeout)
		defer cancel()

		res, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))
		if ok, _ := retryable.DefaultRetryPolicy(ctx, res, err); ok {
			err = retry.RetryableError(err)
		}

		return res, err
	}

	return retry.DoValue(req.Context(), r.backoff, operation)
}
