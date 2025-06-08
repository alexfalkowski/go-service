package retry

import (
	"context"

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
	return &RoundTripper{cfg: cfg, RoundTripper: hrt}
}

// RoundTripper for retry.
type RoundTripper struct {
	cfg *Config
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	timeout := time.MustParseDuration(r.cfg.Timeout)
	back := retry.WithMaxRetries(r.cfg.Attempts, retry.NewConstant(time.MustParseDuration(r.cfg.Backoff)))
	operation := func(ctx context.Context) (*http.Response, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		res, err := r.RoundTripper.RoundTrip(req.WithContext(ctx))
		if ok, _ := retryable.DefaultRetryPolicy(ctx, res, err); ok {
			err = retry.RetryableError(err)
		}

		return res, err
	}

	return retry.DoValue(req.Context(), back, operation)
}
