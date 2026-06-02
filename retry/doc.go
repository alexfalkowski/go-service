// Package retry provides shared retry helpers for go-service.
//
// It wraps retry/backoff helpers such as [Backoff], [NewConstant],
// [WithMaxRetries], [Do], and [DoValue] behind the go-service import path so
// transport packages can share retry semantics without importing third-party
// retry primitives directly.
package retry
