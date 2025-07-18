package health

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// ServerParams for health.
type ServerParams struct {
	di.In
	Observer *Observer `optional:"true"`
}

// NewServer creates a new gRPC health server.
func NewServer(params ServerParams) *Server {
	if params.Observer == nil {
		return nil
	}

	return &Server{observer: params.Observer}
}

// Server represents a gRPC health server.
type Server struct {
	health.UnimplementedHealthServer
	observer *Observer
}

// Check the health.
func (s *Server) Check(_ context.Context, _ *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	var status health.HealthCheckResponse_ServingStatus
	if err := s.observer.Error(); err != nil {
		status = health.HealthCheckResponse_NOT_SERVING
	} else {
		status = health.HealthCheckResponse_SERVING
	}

	return &health.HealthCheckResponse{Status: status}, nil
}

// Watch the health.
func (s *Server) Watch(_ *health.HealthCheckRequest, w health.Health_WatchServer) error {
	var status health.HealthCheckResponse_ServingStatus
	if err := s.observer.Error(); err != nil {
		status = health.HealthCheckResponse_NOT_SERVING
	} else {
		status = health.HealthCheckResponse_SERVING
	}

	return w.Send(&health.HealthCheckResponse{Status: status})
}
