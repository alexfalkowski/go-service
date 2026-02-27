package server

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
)

// Server defines the minimal interface required by Service to manage a transport server.
//
// Implementations are typically thin adapters over concrete servers such as `net/http.Server` or `grpc.Server`.
// The embedded fmt.Stringer is expected to return a human-readable address or identifier (used for logging).
type Server interface {
	fmt.Stringer

	// Serve starts serving requests.
	//
	// Serve is expected to block until the server stops. It should return a non-nil error when serving
	// terminates unexpectedly (for example due to a listener failure). Graceful shutdown should typically
	// cause Serve to return nil or a well-understood sentinel error, depending on the underlying server.
	Serve() error

	// Shutdown stops the server gracefully.
	//
	// Shutdown should attempt to stop the server without abruptly dropping active requests, respecting ctx
	// for deadlines/cancellation. It should return any error encountered during shutdown.
	Shutdown(ctx context.Context) error
}
