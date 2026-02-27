package health

import (
	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-health/v2/subscriber"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// ServerParams defines dependencies for constructing the gRPC health `Server` implementation.
//
// It is an Fx parameter struct (`di.In`) that provides the underlying `*health.Server`
// from `github.com/alexfalkowski/go-health/v2/server`, which maintains health observers.
type ServerParams struct {
	di.In

	// Server is the underlying health server that stores and exposes health observers.
	//
	// It is expected to be non-nil when this gRPC health service is wired.
	Server *health.Server
}

// NewServer constructs a new gRPC health `Server` implementation.
//
// The returned server implements the standard gRPC health protocol service
// (`grpc.health.v1.Health`) and delegates health state lookups to the provided
// underlying `*health.Server`.
func NewServer(params ServerParams) *Server {
	return &Server{server: params.Server}
}

// Server implements the standard gRPC health protocol service.
//
// It exposes health state by querying observers registered in the underlying `*health.Server`.
// The service is typically used by load balancers and orchestration systems to determine whether
// a server is serving traffic.
type Server struct {
	v1.UnimplementedHealthServer
	server *health.Server
}

// Check returns the health status for a single service.
//
// The requested service name is taken from `req.GetService()`. The health state is resolved by looking up
// an observer from the underlying `*health.Server` using the "grpc" transport kind.
//
// Error mapping:
//   - If the requested service does not exist, it returns `codes.NotFound`.
//   - Otherwise, it returns a `HealthCheckResponse` whose status is derived from the observer's error state.
func (s *Server) Check(_ context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	observer, err := s.server.Observer(req.GetService(), "grpc")
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &v1.HealthCheckResponse{Status: s.status(observer)}, nil
}

// List returns the health status for all registered services.
//
// It enumerates all observers registered for the "grpc" transport kind and returns their current serving
// statuses.
func (s *Server) List(_ context.Context, req *v1.HealthListRequest) (*v1.HealthListResponse, error) {
	res := &v1.HealthListResponse{Statuses: map[string]*v1.HealthCheckResponse{}}
	for name, observer := range s.server.Observers("grpc") {
		res.Statuses[name] = &v1.HealthCheckResponse{Status: s.status(observer)}
	}

	return res, nil
}

// Watch streams the current health status for a single service.
//
// Note: this implementation sends a single response with the current status and then returns.
// It does not continuously stream updates over time.
func (s *Server) Watch(req *v1.HealthCheckRequest, w v1.Health_WatchServer) error {
	observer, err := s.server.Observer(req.GetService(), "grpc")
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	return w.Send(&v1.HealthCheckResponse{Status: s.status(observer)})
}

func (s *Server) status(observer *subscriber.Observer) v1.HealthCheckResponse_ServingStatus {
	if err := observer.Error(); err != nil {
		return v1.HealthCheckResponse_NOT_SERVING
	}

	return v1.HealthCheckResponse_SERVING
}
