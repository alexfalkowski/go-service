package retry

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	config "github.com/alexfalkowski/go-service/v2/transport/retry"
	retryable "github.com/hashicorp/go-retryablehttp"
)

// Config is an alias for [github.com/alexfalkowski/go-service/v2/transport/retry.Config].
//
// It describes the retry policy used by NewRoundTripper:
//   - `Attempts`: maximum number of attempts including the initial attempt.
//   - `Backoff`: base delay between retries.
type Config = config.Config

// Policy decides whether req is eligible for retry.
//
// Policies describe operation safety, not transient failure classification. The retry transport still only retries
// retryable responses/errors after the policy allows the logical request.
// When multiple non-nil policies are provided, all must allow the request. Nil policies are ignored; if every
// provided policy is nil, the default [IdempotentRequests] policy is used.
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

// NewRoundTripper constructs a RoundTripper that applies retries.
//
// The constructed RoundTripper wraps hrt and, for each request:
//   - checks whether the request is eligible for retry using policies, and
//   - retries responses and status errors with retryable HTTP status codes, and
//   - retries recoverable transport errors using `retryablehttp.DefaultRetryPolicy`, and
//   - waits a jittered backoff derived from `cfg.GetBackoff()` between attempts.
//
// Attempts/backoff:
// `cfg.Attempts` is interpreted as the total attempt count (initial attempt + retries). Since
// the shared retry helper expects a retry count, NewRoundTripper converts it via `cfg.MaxRetries()`
// before wrapping the jittered backoff in `WithMaxRetries`.
//
// Policy behavior:
// When no policy is provided, only side-effect-safe requests are eligible for retry: safe HTTP methods, or
// requests carrying a request-id idempotency contract. Callers that need different behavior can pass an
// explicit policy.
// Multiple non-nil policies are composed with logical AND, so any denying policy makes the request
// non-retryable. Nil policies are ignored, and the default policy is used only when no non-nil policy is
// supplied.
//
// Exhaustion behavior:
//   - If retries are exhausted after retryable HTTP responses, the final retryable response is returned.
//   - If retries are exhausted after retryable transport errors, the original transport error is returned.
//
// Retry-After behavior:
// Valid Retry-After seconds values or HTTP-date values greater than the minimum jittered backoff suppress the
// retry and return the current response to the caller. Invalid, absent, elapsed, or smaller Retry-After values
// do not suppress a retry.
func NewRoundTripper(cfg *Config, hrt http.RoundTripper, policies ...Policy) *RoundTripper {
	return &RoundTripper{
		RoundTripper: hrt,
		backoff:      cfg.GetBackoff(),
		policy:       composePolicy(policies),
		maxRetries:   cfg.MaxRetries(),
	}
}

// RoundTripper wraps an underlying [http.RoundTripper] and retries requests according to its configuration.
//
// Each attempt:
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
	maxRetries uint64
}

// RoundTrip executes the request and retries according to the configured backoff policy.
//
// For each attempt:
//   - it clones the request and replays the body when needed,
//   - it checks whether the response/status error carries a retryable HTTP status code, and
//   - it uses `retryablehttp.DefaultRetryPolicy` for transport error classification, and
//   - if retryable, it returns a retryable error so the backoff schedules another attempt.
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

	res, err, retrying := r.attempt(req.Context(), req, attempt)
	if !retrying {
		return res, err, false
	}

	res, err = r.retry(req.Context(), req, attempt, err)
	return res, err, false
}

func (r *RoundTripper) retry(ctx context.Context, req *http.Request, attempt *roundTripAttempt, retryErr error) (*http.Response, error) {
	firstRetry := true
	operation := func(ctx context.Context) (*http.Response, error) {
		if firstRetry {
			firstRetry = false
			return nil, retry.RetryableError(retryErr)
		}

		res, err, retrying := r.attempt(ctx, req, attempt)
		if retrying {
			return nil, retry.RetryableError(err)
		}

		return res, err
	}

	backoff := retry.WithJitterPercent(config.DefaultJitterPercent, retry.NewConstant(r.backoff))
	backoff = retry.WithMaxRetries(r.maxRetries, backoff)
	return retry.DoValue(ctx, backoff, operation)
}

func (r *RoundTripper) attempt(ctx context.Context, req *http.Request, attempt *roundTripAttempt) (*http.Response, error, bool) {
	attemptReq, err := request(req, ctx, attempt.attempt)
	if err != nil {
		return nil, err, false
	}

	res, err := r.RoundTripper.RoundTrip(attemptReq)
	if retryErr := attempt.retry(ctx, req, res, err, r.backoff, r.maxRetries); retryErr != nil {
		return nil, retryErr, true
	}

	return res, err, false
}

type roundTripAttempt struct {
	attempt uint64
}

func (a *roundTripAttempt) retry(ctx context.Context, req *http.Request, res *http.Response, err error, backoff time.Duration, maxRetries uint64) error {
	ok, retryErr := shouldRetryAttempt(ctx, res, err)
	if !ok {
		return nil
	}
	if !canRetry(req) {
		return nil
	}

	retryErr = statusError(retryErr, err)
	if res == nil {
		a.attempt++
		return retryErr
	}
	if a.attempt >= maxRetries {
		return nil
	}
	if retryAfterDelayExceedsBackoff(res, backoff) {
		return nil
	}

	closeResponse(res)
	a.attempt++
	return retryErr
}

func request(req *http.Request, ctx context.Context, attempt uint64) (*http.Request, error) {
	if attempt == 0 {
		return req.WithContext(ctx), nil
	}

	clonedReq := req.Clone(ctx)
	if req.GetBody == nil {
		return clonedReq, nil
	}

	retryBody, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	clonedReq.Body = retryBody
	return clonedReq, nil
}

func shouldRetryAttempt(ctx context.Context, res *http.Response, err error) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if errors.Is(err, http.ErrUseLastResponse) {
		return false, err
	}

	if isTransportError(res, err) {
		return retryable.DefaultRetryPolicy(ctx, res, err)
	}

	if err != nil {
		if status.IsLocalError(err) {
			return false, err
		}

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
	if res != nil {
		body.Close(res.Body)
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
