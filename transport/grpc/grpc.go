package grpc

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/net/grpc"
)

// DialOption is an alias for net/grpc.DialOption.
//
// It is exposed so callers can use raw dial options through the transport/grpc
// import path when they do not need the higher-level ClientOption stack.
type DialOption = grpc.DialOption

// NewClientConn constructs a low-level gRPC client channel using raw dial options.
//
// This forwards to net/grpc.NewClient. No I/O is performed during construction;
// the returned ClientConn connects automatically when it is used for RPCs (or
// when Connect is called explicitly).
//
// It is intended for callers that want the transport/grpc import path without
// opting into the higher-level ClientOption stack used by NewClient.
func NewClientConn(target string, opts ...DialOption) (*ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

// NewInsecureCredentials returns transport credentials that disable transport security.
//
// This forwards to net/grpc.NewInsecureCredentials and is intended for local
// development, tests, or deployments where transport security is handled
// out-of-band.
func NewInsecureCredentials() grpc.TransportCredentials {
	return grpc.NewInsecureCredentials()
}

// NewTLS constructs TLS transport credentials from c.
//
// This forwards to net/grpc.NewTLS.
func NewTLS(c *tls.Config) grpc.TransportCredentials {
	return grpc.NewTLS(c)
}

// WithTransportCredentials returns a DialOption that configures client-side transport credentials.
//
// This forwards to net/grpc.WithTransportCredentials. For TLS, use NewTLS to
// construct the credentials from a *tls.Config.
func WithTransportCredentials(creds grpc.TransportCredentials) DialOption {
	return grpc.WithTransportCredentials(creds)
}
