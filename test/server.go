package test

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/pkg/security/meta"
)

// ErrInvalidToken ...
var ErrInvalidToken = errors.New("invalid token")

// NewServer ...
func NewServer(verifyAuth bool) *Server {
	return &Server{verifyAuth: verifyAuth}
}

// Server ...
type Server struct {
	verifyAuth bool
	UnimplementedGreeterServer
}

// SayHello ...
func (s *Server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	if s.verifyAuth && meta.AuthorizedParty(ctx) == "" {
		return nil, ErrInvalidToken
	}

	return &HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())}, nil
}

// SayStreamHello ...
func (s *Server) SayStreamHello(stream Greeter_SayStreamHelloServer) error {
	if s.verifyAuth && meta.AuthorizedParty(stream.Context()) == "" {
		return ErrInvalidToken
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())})
}
