package context

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/time"
)

type (
	// CancelFunc is an alias for context.CancelFunc.
	//
	// It cancels a Context created by WithDeadline or WithTimeout, releasing resources associated with it.
	// As with the standard library, you should call the returned CancelFunc as soon as the operations
	// running in the derived context complete.
	CancelFunc = context.CancelFunc

	// Context is an alias for context.Context.
	//
	// It carries deadlines, cancellation signals, and request-scoped values across API boundaries.
	// The semantics are identical to the standard library context package.
	Context = context.Context

	// Key is a typed helper for storing values in a context.
	//
	// Using a distinct key type reduces accidental collisions when multiple packages store values in the
	// same context. Prefer defining keys as unexported package variables, e.g.:
	//
	//	var userIDKey context.Key = "user_id"
	//
	// and retrieving values with type assertions at the call site.
	Key string
)

// Canceled is an alias for context.Canceled.
//
// It is returned by Context.Err when the context is canceled.
//
//nolint:errname
var Canceled = context.Canceled

// Background returns an empty, non-cancelable context.
//
// This is a thin wrapper around context.Background. Use it as the top-level parent for main,
// initialization, and tests. Request handlers should typically use the request-provided context instead.
func Background() Context {
	return context.Background()
}

// WithDeadline returns a copy of parent with the deadline adjusted to be no later than d.
//
// This is a thin wrapper around context.WithDeadline. The returned CancelFunc should be called to release
// resources associated with the derived context.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return context.WithDeadline(parent, d)
}

// WithTimeout returns a copy of parent with a timeout applied.
//
// This is a thin wrapper around context.WithTimeout. The returned CancelFunc should be called to release
// resources associated with the derived context.
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithValue returns a copy of parent in which the value associated with key is val.
//
// This is a thin wrapper around context.WithValue.
//
// Best practice: use context values only for request-scoped data that transits process boundaries
// (like trace IDs). Do not use them to pass optional parameters to functions.
func WithValue(parent Context, key, val any) Context {
	return context.WithValue(parent, key, val)
}
