// Package retry provides HTTP retry middleware for go-service clients.
//
// The package exposes a transport-level retrying `http.RoundTripper` that:
//   - applies a per-attempt timeout,
//   - retries responses and HTTP status errors only for selected HTTP status codes,
//   - keeps `retryablehttp.DefaultRetryPolicy` classification for transport errors,
//   - replays request bodies on subsequent attempts when `req.GetBody` is available, and
//   - preserves the original caller-visible response or error shape when retries are exhausted.
//
// In particular, when retries are triggered by HTTP responses (for example 429/503), the wrapper keeps the first
// retryable response and returns that response if all retries fail. When retries are triggered by transport errors,
// the original error is preserved.
package retry
