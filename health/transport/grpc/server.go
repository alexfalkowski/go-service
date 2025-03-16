package grpc

import (
	"context"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// NewServer creates a new gRPC health server.
func NewServer(ob *Observer) *Server {
	if ob == nil {
		return nil
	}

	return &Server{ob: ob}
}

// Server represents a gRPC health server.
type Server struct {
	ob *Observer
}

// Check the health.
func (s *Server) Check(_ context.Context, _ *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	status := &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}

	if err := s.ob.Error(); err != nil {
		status.Status = health.HealthCheckResponse_NOT_SERVING
	}

	return status, nil
}

// Watch the health.
func (s *Server) Watch(_ *health.HealthCheckRequest, w health.Health_WatchServer) error {
	status := &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}

	if err := s.ob.Error(); err != nil {
		status.Status = health.HealthCheckResponse_NOT_SERVING
	}

	return w.Send(status)
}
