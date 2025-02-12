package test

import (
	"context"

	"github.com/alexfalkowski/go-service/time"
)

// Timeout for tests.
func Timeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Minute)
}
