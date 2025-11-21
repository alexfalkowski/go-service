package context

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/v2/time"
)

type (
	// CancelFunc is an alias for context.CancelFunc.
	CancelFunc = context.CancelFunc

	// Context is an alias for context.Context.
	Context = context.Context

	// Key is used to store items in context.
	Key string
)

// IsCanceledError returns true if the error is a context canceled error.
func IsCanceledError(err error) bool {
	return errors.Is(err, context.Canceled)
}

// Background is an alias for context.WithValue.
func Background() Context {
	return context.Background()
}

// WithDeadline is an alias for context.WithDeadline.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return context.WithDeadline(parent, d)
}

// WithTimeout is an alias for context.WithTimeout.
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithValue is an alias for context.WithValue.
func WithValue(parent Context, key, val any) Context {
	return context.WithValue(parent, key, val)
}
