package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
)

// DefaultTimeout is the standard timeout used by test contexts in this package.
const DefaultTimeout = 10 * time.Minute

// ErrTimeout is the cause recorded when a test helper context times out.
var ErrTimeout = errors.New("test: timeout")

// Timeout returns a child context with DefaultTimeout applied.
func Timeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(parent, DefaultTimeout, ErrTimeout)
}
