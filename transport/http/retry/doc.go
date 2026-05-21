// Package retry provides HTTP retry middleware for go-service clients.
//
// The package exposes a transport-level retrying `http.RoundTripper` that:
//   - applies a per-attempt timeout,
//   - retries responses and HTTP status errors only for selected HTTP status codes,
//   - keeps `retryablehttp.DefaultRetryPolicy` classification for transport errors,
//   - replays request bodies on subsequent attempts when `req.GetBody` is available, and
//   - returns a caller-visible response or error when retries are exhausted.
//
// In particular, when retries are triggered by HTTP responses (for example 429/503), the wrapper closes
// intermediate retryable responses before scheduling the next attempt and returns the final retryable response if
// all attempts fail. When retries are triggered by transport errors, the original error is preserved.
//
// Backward compatibility: if no policy is passed to NewRoundTripper, all requests are eligible for retry when
// they hit a retryable transport error or status. New callers that only want side-effect-safe retries should pass
// IdempotentRequests, SafeMethods, or another explicit policy.
package retry
