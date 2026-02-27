package health

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterParams defines dependencies for registering the gRPC health service.
//
// It is an Fx parameter struct (`di.In`) used to wire the standard gRPC health protocol service
// (`grpc.health.v1.Health`) into a server.
//
// Fields may be nil depending on configuration/wiring; `Register` is a no-op unless both are provided.
type RegisterParams struct {
	di.In

	// Registrar is the gRPC service registrar used to register the health service implementation.
	//
	// In practice this is typically the transport gRPC server (or something wrapping it) that implements
	// `grpc.ServiceRegistrar`.
	Registrar grpc.ServiceRegistrar

	// Server is the health service implementation to register.
	Server *Server
}

// Register registers the gRPC health service with the provided registrar.
//
// If either `params.Registrar` or `params.Server` is nil, Register does nothing. This allows health wiring
// to be optional in DI graphs where the gRPC server may be disabled.
func Register(params RegisterParams) {
	if params.Registrar == nil || params.Server == nil {
		return
	}

	health.RegisterHealthServer(params.Registrar, params.Server)
}
