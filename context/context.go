package context

import "context"

type (
	// CancelFunc is an alias for context.CancelFunc.
	CancelFunc = context.CancelFunc

	// Context is an alias for context.Context.
	Context = context.Context

	// Key is used to store items in context.
	Key string
)

var (
	// Background is an alias for context.WithValue.
	Background = context.Background

	// WithDeadline is an alias for context.WithDeadline.
	WithDeadline = context.WithDeadline

	// WithTimeout is an alias for context.WithTimeout.
	WithTimeout = context.WithTimeout

	// WithValue is an alias for context.WithValue.
	WithValue = context.WithValue
)
