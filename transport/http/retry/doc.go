// Package retry provides HTTP retry middleware and wiring for go-service.
//
// This package implements an http.RoundTripper that retries requests according to a retry configuration
// (for example timeout, backoff, and attempt count) and is intended to be composed into HTTP clients.
package retry
