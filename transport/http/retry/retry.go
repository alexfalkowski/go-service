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
// It describes the retry policy used by NewRoundTripper:
//   - `Attempts`: maximum number of retries after the initial attempt.
//   - `Timeout`: per-attempt timeout duration string.
//   - `Backoff`: constant delay between retries.
type Config = config.Config

// ErrInvalidStatusCode indicates an HTTP response status code that is considered invalid in the context
// of a retry decision.
//
// This sentinel is used when the retry policy decides a response should be retried but does not provide
// a more specific error value. In practice this happens for retryable HTTP response codes such as
// `429 Too Many Requests` or retryable `5xx` responses.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

// NewRoundTripper constructs a RoundTripper that applies per-attempt timeouts and retries.
//
// The constructed RoundTripper wraps hrt and, for each request:
//   - applies a per-attempt timeout derived from `cfg.Timeout`, and
//   - retries when `retryablehttp.DefaultRetryPolicy` deems the response/error retryable, and
//   - waits a constant backoff derived from `cfg.Backoff` between attempts.
//
// Attempts/backoff:
// `cfg.Attempts` configures the maximum number of retries performed by the backoff policy (not counting
// the initial attempt). Backoff is implemented using `sethvargo/go-retry` with a constant backoff wrapped
// in `WithMaxRetries`.
//
// Exhaustion behavior:
//   - If retries are exhausted after retryable HTTP responses, the first retryable response is returned.
//   - If retries are exhausted after retryable transport errors, the original transport error is returned.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper) *RoundTripper {
	timeout := time.MustParseDuration(cfg.Timeout)
	backoff := retry.WithMaxRetries(cfg.Attempts, retry.NewConstant(time.MustParseDuration(cfg.Backoff)))

	return &RoundTripper{RoundTripper: hrt, timeout: timeout, backoff: backoff}
}

// RoundTripper wraps an underlying `http.RoundTripper` and retries requests according to its configuration.
//
// Each attempt:
//   - runs the underlying `RoundTrip` with a derived per-attempt timeout,
//   - re-creates the request body for subsequent attempts when `req.GetBody` is available, and
//   - asks `retryablehttp.DefaultRetryPolicy` whether the result should be retried.
//
// Important:
// Many HTTP requests are not safe to retry unless they are idempotent or the server supports safe retries.
// This transport does not attempt to determine idempotency; it only applies the configured retry policy.
type RoundTripper struct {
	http.RoundTripper
	backoff retry.Backoff
	timeout time.Duration
}

// RoundTrip executes the request and retries according to the configured backoff policy.
//
// For each attempt:
//   - it derives a child context with a timeout (`r.timeout`) and executes the underlying RoundTripper with it,
//   - it clones the request and replays the body when needed,
//   - it asks `retryablehttp.DefaultRetryPolicy` whether the response/error is retryable, and
//   - if retryable, it returns a retryable error (`retry.RetryableError`) so the backoff schedules another attempt.
//
// When the policy deems the result non-retryable, RoundTrip returns the response/error as produced by the
// underlying transport.
//
// When retries are exhausted:
//   - the first retryable response is returned for response-based failures, preserving the original status/body,
//   - the original transport error is returned for transport-level failures.
//
// Callers should ensure that request bodies are replayable when retries are enabled. For subsequent attempts
// this implementation relies on `req.GetBody`; when it is nil for a request with a body, the first attempt
// is still executed but any retryable result is treated as non-retryable to avoid reusing a consumed body.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	attempt := &roundTripAttempt{}

	operation := func(ctx context.Context) (*http.Response, error) {
		return r.attempt(req, ctx, attempt)
	}

	res, err := retry.DoValue(req.Context(), r.backoff, operation)
	return attempt.finalize(res, err)
}

func (r *RoundTripper) attempt(req *http.Request, ctx context.Context, attempt *roundTripAttempt) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	attemptReq, err := attempt.request(req, ctx)
	if err != nil {
		return nil, err
	}

	res, err := r.RoundTripper.RoundTrip(attemptReq)
	if retryErr := attempt.retry(ctx, req, res, err); retryErr != nil {
		return nil, retryErr
	}

	return res, err
}

func request(req *http.Request, ctx context.Context, attempt int) (*http.Request, error) {
	if attempt == 0 && req.GetBody == nil {
		return req.WithContext(ctx), nil
	}

	cloned := req.Clone(ctx)
	if req.GetBody == nil {
		return cloned, nil
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	cloned.Body = body
	return cloned, nil
}

type roundTripAttempt struct {
	first   *http.Response
	attempt int
}

func (a *roundTripAttempt) request(req *http.Request, ctx context.Context) (*http.Request, error) {
	attemptReq, err := request(req, ctx, a.attempt)
	if err != nil {
		return nil, err
	}

	return attemptReq, nil
}

func (a *roundTripAttempt) retry(ctx context.Context, req *http.Request, res *http.Response, err error) error {
	ok, retryErr := retryable.DefaultRetryPolicy(ctx, res, err)
	if !ok {
		return nil
	}
	if !canRetry(req) {
		return nil
	}

	a.attempt++
	retryErr = statusError(retryErr, err)
	if res == nil {
		return retry.RetryableError(retryErr)
	}

	a.keepFirst(res)
	return retry.RetryableError(responseError{resp: a.first, err: retryErr})
}

func canRetry(req *http.Request) bool {
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
}

func (a *roundTripAttempt) finalize(res *http.Response, err error) (*http.Response, error) {
	if err == nil {
		a.closeFirst(res)
		return res, nil
	}

	if re, ok := errors.AsType[responseError](err); ok {
		return re.resp, nil
	}

	closeResponse(a.first)
	return res, err
}

func (a *roundTripAttempt) keepFirst(res *http.Response) {
	if a.first == nil {
		a.first = res
		return
	}

	closeResponse(res)
}

func (a *roundTripAttempt) closeFirst(res *http.Response) {
	if a.first == nil || a.first == res {
		return
	}

	closeResponse(a.first)
}

func statusError(retryErr, err error) error {
	if retryErr != nil {
		return retryErr
	}

	if err != nil {
		return err
	}

	return ErrInvalidStatusCode
}

func closeResponse(res *http.Response) {
	if res != nil && res.Body != nil {
		_ = res.Body.Close()
	}
}
