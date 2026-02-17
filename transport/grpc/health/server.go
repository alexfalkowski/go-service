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

// ServerParams defines dependencies for constructing the gRPC health server.
type ServerParams struct {
	di.In
	Server *health.Server
}

// NewServer creates a new gRPC health server.
func NewServer(params ServerParams) *Server {
	return &Server{server: params.Server}
}

// Server represents a gRPC health server.
type Server struct {
	v1.UnimplementedHealthServer
	server *health.Server
}

// Check returns the health status for a single service.
func (s *Server) Check(_ context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	observer, err := s.server.Observer(req.GetService(), "grpc")
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &v1.HealthCheckResponse{Status: s.status(observer)}, nil
}

// List returns the health status for all registered services.
func (s *Server) List(_ context.Context, req *v1.HealthListRequest) (*v1.HealthListResponse, error) {
	res := &v1.HealthListResponse{Statuses: map[string]*v1.HealthCheckResponse{}}
	for name, observer := range s.server.Observers("grpc") {
		res.Statuses[name] = &v1.HealthCheckResponse{Status: s.status(observer)}
	}

	return res, nil
}

// Watch the health.
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
