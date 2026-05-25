// Package retry provides shared retry helpers for go-service.
//
// It wraps the repository's retry/backoff implementation behind the go-service
// import path so transport packages can share retry semantics without importing
// third-party retry primitives directly.
package retry
