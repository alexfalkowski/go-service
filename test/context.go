package test

import (
	"context"
	"time"
)

// Timeout for tests.
func Timeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Minute)
}
