package retry

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	"github.com/alexfalkowski/go-sync"
	retryable "github.com/hashicorp/go-retryablehttp"
	"github.com/sethvargo/go-retry"
)

// Config is an alias for `github.com/alexfalkowski/go-service/v2/transport/retry.Config`.
//
// It describes the retry policy used by NewRoundTripper:
//   - `Attempts`: maximum number of attempts including the initial attempt.
//   - `Timeout`: per-attempt timeout duration.
//   - `Backoff`: constant delay between retries.
type Config = config.Config

// Policy decides whether req is eligible for retry.
//
// Policies describe operation safety, not transient failure classification. The retry transport still only retries
// retryable responses/errors after the policy allows the logical request.
type Policy func(req *http.Request) bool

// SafeMethods allows retries for HTTP methods that should not have request side effects.
func SafeMethods(req *http.Request) bool {
	switch req.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// HasRequestID allows retries when request metadata contains a request id.
//
// In go-service, Request-Id identifies the logical request, not an individual wire attempt. The client metadata
// middleware installs it outside the retry middleware, so all retry attempts for one logical request share the same
// value. Services that retry writes should treat Request-Id as the idempotency key and deduplicate repeated attempts
// by it.
func HasRequestID(req *http.Request) bool {
	return !meta.RequestID(req.Context()).IsEmpty()
}

// IdempotentRequests allows safe methods, or requests with a request-id idempotency contract.
func IdempotentRequests(req *http.Request) bool {
	return SafeMethods(req) || HasRequestID(req)
}

// ErrInvalidStatusCode indicates an HTTP response status code that is considered invalid in the context
// of a retry decision.
//
// This sentinel is used when the retry policy decides a response should be retried but does not provide
// a more specific error value. In practice this happens for retryable HTTP response codes such as
// `429 Too Many Requests` or `503 Service Unavailable`.
var ErrInvalidStatusCode = errors.New("retry: invalid status code")

// ErrAttemptTimeout is the cause recorded when a retry attempt times out.
var ErrAttemptTimeout = fmt.Errorf("retry: attempt timeout: %w", sync.ErrTimeout)

// NewRoundTripper constructs a RoundTripper that applies per-attempt timeouts and retries.
//
// The constructed RoundTripper wraps hrt and, for each request:
//   - checks whether the request is eligible for retry using policies, and
//   - applies a per-attempt timeout derived from `cfg.GetTimeout()`, and
//   - retries responses and status errors with retryable HTTP status codes, and
//   - retries recoverable transport errors using `retryablehttp.DefaultRetryPolicy`, and
//   - waits a constant backoff derived from `cfg.GetBackoff()` between attempts.
//
// Attempts/backoff:
// `cfg.Attempts` is interpreted as the total attempt count (initial attempt + retries). Since
// `sethvargo/go-retry` expects a retry count, NewRoundTripper converts it via `cfg.MaxRetries()`
// before wrapping the constant backoff in `WithMaxRetries`.
//
// Policy behavior:
// When no policy is provided, only side-effect-safe requests are eligible for retry: safe HTTP methods, or
// requests carrying a request-id idempotency contract. Callers that need different behavior can pass an
// explicit policy.
//
// Exhaustion behavior:
//   - If retries are exhausted after retryable HTTP responses, the final retryable response is returned.
//   - If retries are exhausted after retryable transport errors, the original transport error is returned.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper, policies ...Policy) *RoundTripper {
	return &RoundTripper{
		RoundTripper: hrt,
		backoff:      cfg.GetBackoff(),
		policy:       composePolicy(policies),
		timeout:      cfg.GetTimeout(),
		maxRetries:   cfg.MaxRetries(),
	}
}

// RoundTripper wraps an underlying `http.RoundTripper` and retries requests according to its configuration.
//
// Each attempt:
//   - runs the underlying `RoundTrip` with a derived per-attempt timeout,
//   - re-creates the request body for subsequent attempts when `req.GetBody` is available, and
//   - retries only selected HTTP status codes for responses and status errors.
//   - retries recoverable transport errors using `retryablehttp.DefaultRetryPolicy`.
//
// Important:
// Many HTTP requests are not safe to retry unless they are idempotent or the server supports safe retries.
// This transport defaults to IdempotentRequests and only applies broader retry behavior when callers provide
// an explicit policy.
type RoundTripper struct {
	http.RoundTripper
	policy     Policy
	backoff    time.Duration
	timeout    time.Duration
	maxRetries uint64
}

// RoundTrip executes the request and retries according to the configured backoff policy.
//
// For each attempt:
//   - it derives a child context with a timeout (`r.timeout`) and executes the underlying RoundTripper with it,
//   - it clones the request and replays the body when needed,
//   - it checks whether the response/status error carries a retryable HTTP status code, and
//   - it uses `retryablehttp.DefaultRetryPolicy` for transport error classification, and
//   - if retryable, it returns a retryable error (`retry.RetryableError`) so the backoff schedules another attempt.
//
// When the policy deems the result non-retryable, RoundTrip returns the response/error as produced by the
// underlying transport.
//
// When retries are exhausted:
//   - the final retryable response is returned for response-based failures,
//   - the original transport error is returned for transport-level failures.
//
// Callers should ensure that request bodies are replayable when retries are enabled. For subsequent attempts
// this implementation relies on `req.GetBody`; when it is nil for a request with a body, the first attempt
// is still executed but any retryable result is treated as non-retryable to avoid reusing a consumed body.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.ClosingRoundTripper(r.roundTrip).RoundTrip(req)
}

func (r *RoundTripper) roundTrip(req *http.Request) (*http.Response, error, bool) {
	if err := req.Context().Err(); err != nil {
		return nil, err, true
	}

	if r.policy != nil && !r.policy(req) {
		res, err := r.RoundTripper.RoundTrip(req)
		return res, err, false
	}

	attempt := &roundTripAttempt{}

	operation := func(ctx context.Context) (*http.Response, error) {
		return r.attempt(ctx, req, attempt)
	}

	backoff := retry.WithMaxRetries(r.maxRetries, retry.NewConstant(r.backoff.Duration()))
	res, err := retry.DoValue(req.Context(), backoff, operation)
	return res, err, false
}

func (r *RoundTripper) attempt(ctx context.Context, req *http.Request, attempt *roundTripAttempt) (*http.Response, error) {
	attemptCtx, cancel := r.withAttemptTimeout(ctx)
	defer cancel()

	attemptReq, err := attempt.request(req, attemptCtx)
	if err != nil {
		return nil, err
	}

	res, err := r.RoundTripper.RoundTrip(attemptReq)
	if retryErr := attempt.retry(ctx, attemptCtx, req, res, err, r.maxRetries); retryErr != nil {
		return nil, retryErr
	}

	return res, err
}

func (r *RoundTripper) withAttemptTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(ctx, r.timeout, ErrAttemptTimeout)
}

func request(req *http.Request, ctx context.Context, attempt uint64) (*http.Request, error) {
	if attempt == 0 {
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
	attempt uint64
}

func (a *roundTripAttempt) request(req *http.Request, ctx context.Context) (*http.Request, error) {
	attemptReq, err := request(req, ctx, a.attempt)
	if err != nil {
		return nil, err
	}

	return attemptReq, nil
}

func (a *roundTripAttempt) retry(ctx, attemptCtx context.Context, req *http.Request, res *http.Response, err error, maxRetries uint64) error {
	ok, retryErr := shouldRetryAttempt(ctx, attemptCtx, res, err)
	if !ok {
		return nil
	}
	if !canRetry(req) {
		return nil
	}

	retryErr = statusError(retryErr, err)
	if res == nil {
		a.attempt++
		return retry.RetryableError(retryErr)
	}
	if a.attempt >= maxRetries {
		return nil
	}

	closeResponse(res)
	a.attempt++
	return retry.RetryableError(retryErr)
}

func shouldRetryAttempt(ctx, attemptCtx context.Context, res *http.Response, err error) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if err := attemptCtx.Err(); err != nil {
		cause := context.Cause(attemptCtx)
		return errors.Is(cause, ErrAttemptTimeout), cause
	}

	if isTransportError(res, err) {
		return retryable.DefaultRetryPolicy(ctx, res, err)
	}

	if err != nil {
		return isRetryableStatusCode(status.Code(err)), err
	}

	if res == nil {
		return false, nil
	}

	return isRetryableStatusCode(res.StatusCode), nil
}

func isRetryableStatusCode(code int) bool {
	return code == http.StatusTooManyRequests || code == http.StatusServiceUnavailable
}

func isTransportError(res *http.Response, err error) bool {
	return res == nil && err != nil && !status.IsError(err)
}

func canRetry(req *http.Request) bool {
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
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

func composePolicy(policies []Policy) Policy {
	filtered := make([]Policy, 0, len(policies))
	for _, policy := range policies {
		if policy != nil {
			filtered = append(filtered, policy)
		}
	}

	if len(filtered) == 0 {
		return IdempotentRequests
	}

	return func(req *http.Request) bool {
		for _, policy := range filtered {
			if !policy(req) {
				return false
			}
		}

		return true
	}
}
