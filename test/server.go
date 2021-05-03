package test

import (
	"context"
	"fmt"
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())}, nil
}
func (s *Server) SayStreamHello(stream Greeter_SayStreamHelloServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())})
}
