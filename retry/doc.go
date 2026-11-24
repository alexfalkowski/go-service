// Package retry provides shared retry helpers for go-service.
//
// It wraps retry/backoff helpers such as [Backoff], [NewBackoff],
// [WithJitterPercent], [WithMaxRetries], [Do], and [DoValue] behind the
// go-service import path so transport packages can share retry semantics
// without importing third-party retry primitives directly.
//
// # Migrating from Backoff.Next
//
// Replace b.Next() with b(attempt), where attempt is the number of retries
// already made. Use b(0) for the first retry wait. Custom backoffs are functions
// of type func(uint64) (time.Duration, bool), using the standard library's
// time.Duration. Each call to [Do] or [DoValue] starts a fresh retry count.
// Use errors.Is to check for the context cause in cancellation errors.
package retry
