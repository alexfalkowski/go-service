package test

import (
	"context"

	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
)

// NewService ...
func NewService() *Service {
	return &Service{}
}

// Service ...
type Service struct {
	v1.UnimplementedGreeterServiceServer
}

// SayHello ...
func (s *Service) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}

// SayStreamHello ...
func (s *Service) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()})
}
