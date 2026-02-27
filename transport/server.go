package transport

import (
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// ServersParams defines dependencies used to collect transport servers.
//
// It is an Fx parameter struct that gathers the transport server constructors that may or may not be enabled
// at runtime (for example, HTTP disabled in config). Each field may be nil depending on configuration.
type ServersParams struct {
	di.In

	// HTTP is the HTTP transport server, if enabled.
	HTTP *http.Server

	// GRPC is the gRPC transport server, if enabled.
	GRPC *grpc.Server

	// Debug is the debug server, if enabled.
	Debug *debug.Server
}

// NewServers collects the enabled transport services.
//
// It returns a slice of `*server.Service` that excludes nil services, so it is safe to pass the returned slice
// to `Register` for lifecycle wiring.
func NewServers(params ServersParams) []*server.Service {
	return slices.AppendNotNil([]*server.Service{}, params.HTTP.GetService(), params.GRPC.GetService(), params.Debug.GetService())
}
