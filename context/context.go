package context

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/time"
)

// CancelCauseFunc is an alias for context.CancelCauseFunc.
//
// It cancels a Context created by WithCancelCause and records the provided cause. Context.Err still
// reports context.Canceled, while Cause exposes the richer diagnostic error.
type CancelCauseFunc = context.CancelCauseFunc

// CancelFunc is an alias for context.CancelFunc.
//
// It cancels a Context created by WithCancel, WithDeadline, or WithTimeout, releasing resources
// associated with it. As with the standard library, you should call the returned CancelFunc as soon
// as the operations running in the derived context complete.
type CancelFunc = context.CancelFunc

// Context is an alias for context.Context.
//
// It carries deadlines, cancellation signals, and request-scoped values across API boundaries.
// The semantics are identical to the standard library context package.
type Context = context.Context

// Key is a typed helper for storing values in a context.
//
// Using a distinct key type reduces accidental collisions when multiple packages store values in the
// same context. Prefer defining keys as unexported package variables, e.g.:
//
//	var userIDKey context.Key = "user_id"
//
// and retrieving values with type assertions at the call site.
type Key string

// Canceled is an alias for context.Canceled.
//
// It is returned by Context.Err when the context is canceled.
//
//nolint:errname
var Canceled = context.Canceled

// DeadlineExceeded is an alias for context.DeadlineExceeded.
//
// It is returned by Context.Err when the context deadline passes. When a cause-aware deadline or timeout
// API is used, Cause can expose a richer diagnostic error while Err still reports DeadlineExceeded.
//
//nolint:errname
var DeadlineExceeded = context.DeadlineExceeded

// Background returns an empty, non-cancelable context.
//
// This is a thin wrapper around context.Background. Use it as the top-level parent for main,
// initialization, and tests. Request handlers should typically use the request-provided context instead.
func Background() Context {
	return context.Background()
}

// Cause returns the reason a context was canceled.
//
// This is a thin wrapper around context.Cause. Cause returns nil until the context is canceled. When a
// cause-aware API is used, Cause can expose a more specific diagnostic error while Err still reports the
// standard cancellation sentinel.
func Cause(ctx Context) error {
	return context.Cause(ctx)
}

// WithCancel returns a derived context that points to the parent context but has a new Done channel.
//
// This is a thin wrapper around context.WithCancel. The returned CancelFunc should be called to release
// resources associated with the derived context.
func WithCancel(parent Context) (Context, CancelFunc) {
	return context.WithCancel(parent)
}

// WithCancelCause returns a derived context that points to the parent context and records an optional cause.
//
// This is a thin wrapper around context.WithCancelCause. Context.Err still reports Canceled, while Cause
// returns the provided diagnostic error (or Canceled when cancel is called with nil).
func WithCancelCause(parent Context) (Context, CancelCauseFunc) {
	return context.WithCancelCause(parent)
}

// WithDeadline returns a copy of parent with the deadline adjusted to be no later than d.
//
// This is a thin wrapper around context.WithDeadline. The returned CancelFunc should be called to release
// resources associated with the derived context.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return context.WithDeadline(parent, d)
}

// WithDeadlineCause returns a copy of parent with the deadline adjusted to be no later than d.
//
// This is a thin wrapper around context.WithDeadlineCause. If the deadline expires, Err reports
// DeadlineExceeded while Cause returns the configured diagnostic error. Calling the returned CancelFunc
// directly cancels the context without setting the configured cause.
func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc) {
	return context.WithDeadlineCause(parent, d, cause)
}

// WithTimeout returns a copy of parent with a timeout applied.
//
// This is a thin wrapper around context.WithTimeout. The returned CancelFunc should be called to release
// resources associated with the derived context.
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return context.WithTimeout(parent, timeout.Duration())
}

// WithTimeoutCause returns a copy of parent with a timeout applied.
//
// This is a thin wrapper around context.WithTimeoutCause. If the timeout expires, Err reports
// DeadlineExceeded while Cause returns the configured diagnostic error. Calling the returned CancelFunc
// directly cancels the context without setting the configured cause.
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc) {
	return context.WithTimeoutCause(parent, timeout.Duration(), cause)
}

// WithoutCancel returns a copy of parent that is not canceled when parent is canceled.
//
// This is a thin wrapper around context.WithoutCancel. The returned context still exposes values from
// parent, but it has no deadline, no Done channel, and Err always reports nil.
func WithoutCancel(parent Context) Context {
	return context.WithoutCancel(parent)
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
