package server

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
)

// Server defines the minimal interface required by Service to manage a transport server.
type Server interface {
	fmt.Stringer

	// Serve starts serving requests.
	Serve() error

	// Shutdown stops the server gracefully.
	Shutdown(ctx context.Context) error
}
