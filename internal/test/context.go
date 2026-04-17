package test

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

// DefaultTimeout is the standard timeout used by test contexts in this package.
const DefaultTimeout = 10 * time.Minute

// ErrTimeout is the cause recorded when a test helper context times out.
var ErrTimeout = fmt.Errorf("test: timeout: %w", sync.ErrTimeout)

// Timeout returns a child context with DefaultTimeout applied.
func Timeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(parent, DefaultTimeout, ErrTimeout)
}
