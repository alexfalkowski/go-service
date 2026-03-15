// Package rand provides cryptographically secure random helpers and wiring for go-service.
//
// This package wraps crypto/rand and provides a Generator that can produce raw
// random bytes and random text tokens.
//
// Use GenerateBytes when callers need binary randomness, and GenerateText when
// callers specifically want text drawn from the package's alphanumeric alphabet.
package rand
