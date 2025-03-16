package server

import (
	"context"
	"fmt"
)

// Serverer allows transports to create their own.
type Serverer interface {
	fmt.Stringer

	// Serve the underlying server.
	Serve() error

	// Shutdown the underlying server.
	Shutdown(ctx context.Context) error
}
