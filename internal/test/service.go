package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
)

// NewService returns the example Greeter gRPC service used by transport tests.
func NewService() *Service {
	return &Service{}
}

// Service is a minimal Greeter service implementation for the generated test protobuf API.
type Service struct {
	v1.UnimplementedGreeterServiceServer
}

// SayHello returns a greeting for the requested name.
func (s *Service) SayHello(ctx context.Context, req *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + req.GetName()}, nil
}

// SayStreamHello reads a single request from the bidi stream and replies with one greeting message.
func (s *Service) SayStreamHello(stream v1.GreeterService_SayStreamHelloServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&v1.SayStreamHelloResponse{Message: "Hello " + req.GetName()})
}
