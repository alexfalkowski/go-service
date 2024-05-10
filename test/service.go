package test

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
)

// ErrInvalidToken ...
var ErrInvalidToken = errors.New("invalid token")

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
	if s.verifyAuth && !meta.IsEqual(Test(ctx), "auth") {
		return nil, ErrInvalidToken
	}

	return &v1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}

// SayStreamHello ...
func (s *Service) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	if s.verifyAuth && !meta.IsEqual(Test(stream.Context()), "auth") {
		return ErrInvalidToken
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()})
}
