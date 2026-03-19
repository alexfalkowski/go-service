package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
)

// DefaultTimeout is the standard timeout used by test contexts in this package.
const DefaultTimeout = 10 * time.Minute

// Timeout returns a background context with DefaultTimeout applied.
func Timeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}
