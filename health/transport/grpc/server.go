package grpc

import (
	"context"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type server struct {
	ob *Observer
}

// Check the health.
func (s *server) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	status := &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}

	if err := s.ob.Error(); err != nil {
		status.Status = health.HealthCheckResponse_NOT_SERVING
	}

	return status, nil
}

// Watch the health.
func (s *server) Watch(req *health.HealthCheckRequest, w health.Health_WatchServer) error {
	status := &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}

	if err := s.ob.Error(); err != nil {
		status.Status = health.HealthCheckResponse_NOT_SERVING
	}

	return w.Send(status)
}
