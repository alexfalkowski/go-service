package server

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
)

// Server allows transports to create their own.
type Server interface {
	fmt.Stringer

	// Serve the underlying server.
	Serve() error

	// Shutdown the underlying server.
	Shutdown(ctx context.Context) error
}
