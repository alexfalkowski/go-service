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

// Config is an alias for `github.com/alexfalkowski/go-service/v2/retry.Config`.
//
// It describes an HTTP retry policy, including:
//   - `Attempts`: maximum number of retries (not total attempts).
//   - `Timeout`: per-attempt timeout duration string.
//   - `Backoff`: backoff duration string used by the backoff strategy.
type Config = config.Config

// ErrInvalidStatusCode indicates an HTTP response status code that is considered invalid in the context
// of a retry decision.
//
// Note: this package relies on `retryablehttp.DefaultRetryPolicy` for retry classification; this error is
// provided for callers that want a stable sentinel value when wiring custom retry logic.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

// NewRoundTripper constructs a RoundTripper that applies per-attempt timeouts and retries.
//
// The constructed RoundTripper wraps hrt and, for each request:
//   - applies a per-attempt timeout derived from `cfg.Timeout`, and
//   - retries when `retryablehttp.DefaultRetryPolicy` deems the response/error retryable,
//     using a constant backoff step derived from `cfg.Backoff`.
//
// Attempts/backoff:
// `cfg.Attempts` configures the maximum number of retries performed by the backoff policy (not counting
// the initial attempt). Backoff is implemented using `sethvargo/go-retry` with a constant backoff wrapped
// in `WithMaxRetries`.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper) *RoundTripper {
	timeout := time.MustParseDuration(cfg.Timeout)
	backoff := retry.WithMaxRetries(cfg.Attempts, retry.NewConstant(time.MustParseDuration(cfg.Backoff)))

	return &RoundTripper{RoundTripper: hrt, timeout: timeout, backoff: backoff}
}

// RoundTripper wraps an underlying `http.RoundTripper` and retries requests according to its configuration.
//
// Each attempt runs the underlying `RoundTrip` with a derived context that has a per-attempt timeout applied.
// Whether an attempt is considered retryable is determined by `retryablehttp.DefaultRetryPolicy`.
//
// Important:
// Many HTTP requests are not safe to retry unless they are idempotent or the server supports safe retries.
// This transport does not attempt to determine idempotency; it only applies the retry policy.
type RoundTripper struct {
	http.RoundTripper
	backoff retry.Backoff
	timeout time.Duration
}

// RoundTrip executes the request and retries according to the configured backoff policy.
//
// For each attempt:
//   - it derives a child context with a timeout (`r.timeout`) and executes the underlying RoundTripper with it,
//   - it asks `retryablehttp.DefaultRetryPolicy` whether the response/error is retryable, and
//   - if retryable, it returns a retryable error (`retry.RetryableError`) so the backoff schedules another attempt.
//
// When the policy deems the result non-retryable, RoundTrip returns the response/error as produced by the
// underlying transport.
//
// Note: since this is implemented at the RoundTripper layer, callers should ensure that request bodies are
// replayable when retries are enabled.
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
