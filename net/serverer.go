package net

import (
	"context"
	"fmt"
)

// Serverer for net.
type Serverer interface {
	fmt.Stringer

	// Serve the underlying server.
	Serve() error

	// Shutdown the underlying server.
	Shutdown(ctx context.Context) error
}
