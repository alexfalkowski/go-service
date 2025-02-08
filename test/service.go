package test

import (
	"context"

	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
)

// NewService ...
func NewService(verifyAuth bool) *Service {
	return &Service{verifyAuth: verifyAuth}
}

// Service ...
type Service struct {
	v1.UnimplementedGreeterServiceServer
	verifyAuth bool
}

// SayHello ...
func (s *Service) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	if s.verifyAuth && Test(ctx).Value() != "auth" {
		return nil, ErrInvalid
	}

	return &v1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}

// SayStreamHello ...
func (s *Service) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	if s.verifyAuth && Test(stream.Context()).Value() != "auth" {
		return ErrInvalid
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()})
}
