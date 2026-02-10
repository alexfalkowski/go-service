package context

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/time"
)

type (
	// CancelFunc is an alias for context.CancelFunc.
	CancelFunc = context.CancelFunc

	// Context is an alias for context.Context.
	Context = context.Context

	// Key is used as a typed key for storing values in a context.
	//
	// This helps avoid key collisions when using context values.
	Key string
)

// Canceled is an alias for context.Canceled.
//
//nolint:errname
var Canceled = context.Canceled

// Background is an alias for context.Background.
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
