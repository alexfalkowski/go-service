package errors

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/redis/go-redis/v9"
)

// ErrExpired is returned when a cache entry exists but is expired.
//
// Drivers may wrap this error; use [IsExpiredError] to classify this condition.
var ErrExpired = errors.New("cache: expired")

// ErrMissing is returned when a cache entry does not exist.
//
// Drivers may wrap this error; use [IsMissingError] to classify this condition.
var ErrMissing = errors.New("cache: missing")

// ErrNotFound is returned when the configured cache driver kind is unknown.
var ErrNotFound = errors.New("cache: driver not found")

// ErrInvalidURL is returned when a cache backend URL cannot be parsed.
var ErrInvalidURL = errors.New("cache: invalid driver url")

// IsExpiredError reports whether err represents an expired cache entry.
//
// This helper exists so higher-level code can treat expired entries as cache misses regardless of the
// underlying backend implementation.
func IsExpiredError(err error) bool {
	return errors.Is(err, ErrExpired)
}

// IsMissingError reports whether err represents a missing cache entry.
//
// This helper normalizes the miss semantics of the backends currently supported by this package,
// including Redis nil replies ([github.com/redis/go-redis/v9.Nil]).
func IsMissingError(err error) bool {
	if errors.Is(err, redis.Nil) {
		return true
	}

	return errors.Is(err, ErrMissing)
}
