// Package retry provides HTTP retry middleware for go-service clients.
//
// The package exposes a transport-level retrying `http.RoundTripper` that:
//   - retries responses and HTTP status errors only for selected HTTP status codes,
//   - keeps `retryablehttp.DefaultRetryPolicy` classification for transport errors,
//   - replays request bodies on subsequent attempts when `req.GetBody` is available, and
//   - returns a caller-visible response or error when retries are exhausted.
//
// In particular, when retries are triggered by HTTP responses (for example 429/503), the wrapper closes
// intermediate retryable responses before scheduling the next attempt and returns the final retryable response if
// all attempts fail. When retries are triggered by transport errors, the original error is preserved.
//
// Retry-After handling: when a retryable HTTP response includes a valid Retry-After seconds value or HTTP-date
// whose delay is greater than the minimum jittered backoff, the response is returned without another attempt.
// Invalid, absent, elapsed, or shorter Retry-After values do not suppress a retry.
//
// Default policy: if no policy is passed to NewRoundTripper, only side-effect-safe requests are eligible for
// retry. This includes safe HTTP methods and requests carrying a Request-Id. In go-service, Request-Id
// identifies the logical request and is stable across retry attempts, so services that retry writes should
// deduplicate by Request-Id. Callers that need different retry eligibility can pass an explicit policy.
package retry
