package grpc

import (
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor.
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type (
	// ClientConn is an alias of grpc.ClientConn.
	ClientConn = grpc.ClientConn

	// DialOption is an alias of grpc.DialOption.
	DialOption = grpc.DialOption

	// EmptyServerOption is an alias of grpc.EmptyServerOption.
	EmptyServerOption = grpc.EmptyServerOption

	// ServerOption is an alias of grpc.ServerOption.
	ServerOption = grpc.ServerOption

	// Server is an alias of grpc.Server.
	Server = grpc.Server

	// ServiceRegistrar is an alias of grpc.ServiceRegistrar.
	ServiceRegistrar = grpc.ServiceRegistrar

	// StreamClientInterceptor is an alias of grpc.StreamClientInterceptor.
	StreamClientInterceptor = grpc.StreamClientInterceptor

	// StreamServerInterceptor is an alias of grpc.StreamServerInterceptor.
	StreamServerInterceptor = grpc.StreamServerInterceptor

	// UnaryClientInterceptor is an alias of grpc.UnaryClientInterceptor.
	UnaryClientInterceptor = grpc.UnaryClientInterceptor

	// UnaryServerInterceptor is an alias of grpc.UnaryServerInterceptor.
	UnaryServerInterceptor = grpc.UnaryServerInterceptor
)

var (
	// ChainUnaryInterceptor is an alias of grpc.ChainUnaryInterceptor.
	ChainUnaryInterceptor = grpc.ChainUnaryInterceptor

	// ChainStreamInterceptor is an alias of grpc.ChainStreamInterceptor.
	ChainStreamInterceptor = grpc.ChainStreamInterceptor

	// Creds is an alias of grpc.Creds.
	Creds = grpc.Creds

	// NewClient is an alias of grpc.NewClient.
	NewClient = grpc.NewClient

	// NewInsecureCredentials is an alias for insecure.NewCredentials.
	NewInsecureCredentials = insecure.NewCredentials

	// NewClient is an alias of credentials.NewTLS.
	NewTLS = credentials.NewTLS

	// TimeoutUnaryClientInterceptor is an alias of timeout.UnaryClientInterceptor.
	TimeoutUnaryClientInterceptor = timeout.UnaryClientInterceptor

	// UseCompressor is an alias of grpc.UseCompressor.
	UseCompressor = grpc.UseCompressor

	// WithChainUnaryInterceptor is an alias of grpc.WithChainUnaryInterceptor.
	WithChainUnaryInterceptor = grpc.WithChainUnaryInterceptor

	// WithChainStreamInterceptor is an alias of grpc.WithChainStreamInterceptor.
	WithChainStreamInterceptor = grpc.WithChainStreamInterceptor

	// WithDefaultCallOptions is an alias of grpc.WithDefaultCallOptions.
	WithDefaultCallOptions = grpc.WithDefaultCallOptions

	// WithUserAgent is an alias of grpc.WithUserAgent.
	WithUserAgent = grpc.WithUserAgent

	// WithTransportCredentials is an alias of grpc.WithTransportCredentials.
	WithTransportCredentials = grpc.WithTransportCredentials
)

// WithKeepaliveParams for grpc.
func WithKeepaliveParams(timeout time.Duration) DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                timeout,
		Timeout:             timeout,
		PermitWithoutStream: true,
	})
}

// NewServer for grpc.
func NewServer(timeout time.Duration, opts ...ServerOption) *Server {
	options := []ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             timeout,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     timeout,
			MaxConnectionAge:      timeout,
			MaxConnectionAgeGrace: timeout,
			Time:                  timeout,
			Timeout:               timeout,
		}),
	}
	options = append(options, opts...)

	server := grpc.NewServer(options...)
	reflection.Register(server)

	return server
}
