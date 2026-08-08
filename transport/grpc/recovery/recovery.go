package recovery

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

// NewServer constructs a gRPC panic recovery interceptor server.
//
// handler converts a recovered panic value into the error returned to the RPC caller. If handler is nil, the
// interceptor returns the upstream default panic error.
func NewServer(handler func(any) error) *Server {
	return &Server{handler: handler}
}

// Server provides gRPC server interceptors that recover panics.
type Server struct {
	handler func(any) error
}

// UnaryInterceptor returns a gRPC unary server interceptor that recovers panics.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(s.handler))
}

// StreamInterceptor returns a gRPC stream server interceptor that recovers panics.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(s.handler))
}
