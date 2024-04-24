package net

import "context"

// Serverer for net.
type Serverer interface {
	// Serve the underlying server.
	Serve() error

	// Shutdown the underlying server.
	Shutdown(ctx context.Context) error
}
