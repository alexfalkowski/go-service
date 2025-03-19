package server

import (
	"context"
	"fmt"
)

// Server allows transports to create their own.
type Server interface {
	fmt.Stringer

	// Serve the underlying server.
	Serve() error

	// Shutdown the underlying server.
	Shutdown(ctx context.Context) error
}
