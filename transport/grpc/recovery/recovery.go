package recovery

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

// NewServer constructs a gRPC panic recovery interceptor server.
//
// handler converts a recovered panic value into the error returned to the RPC caller and must be non-nil.
// UnaryInterceptor and StreamInterceptor always call the upstream
// [github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery.WithRecoveryHandler] with handler,
// which wraps handler in a non-nil callback regardless of whether handler itself is nil. If handler is nil, that
// callback invokes a nil func value when a panic is recovered, which panics again with a nil pointer
// dereference; that second panic happens outside the interceptor's own recover and propagates uncaught, so a
// nil handler crashes the process on any recovered panic instead of returning a safe error. Return a
// client-safe error from handler, for example via [github.com/alexfalkowski/go-service/v2/net/grpc/status.SafeError],
// since gRPC renders an unwrapped handler error as the client-visible status message.
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
