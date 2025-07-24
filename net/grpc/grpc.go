package grpc

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor.
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type (
	// CallOption is an alias of grpc.CallOption.
	CallOption = grpc.CallOption

	// ClientConn is an alias of grpc.ClientConn.
	ClientConn = grpc.ClientConn

	// ClientStream is an alias of grpc.ClientStream.
	ClientStream = grpc.ClientStream

	// DialOption is an alias of grpc.DialOption.
	DialOption = grpc.DialOption

	// EmptyServerOption is an alias of grpc.EmptyServerOption.
	EmptyServerOption = grpc.EmptyServerOption

	// UnaryInvoker is an alias of grpc.UnaryInvoker.
	UnaryInvoker = grpc.UnaryInvoker

	// ServerOption is an alias of grpc.ServerOption.
	ServerOption = grpc.ServerOption

	// Server is an alias of grpc.Server.
	Server = grpc.Server

	// ServerStream is an alias of grpc.ServerStream.
	ServerStream = grpc.ServerStream

	// ServiceRegistrar is an alias of grpc.ServiceRegistrar.
	ServiceRegistrar = grpc.ServiceRegistrar

	// StreamClientInterceptor is an alias of grpc.StreamClientInterceptor.
	StreamClientInterceptor = grpc.StreamClientInterceptor

	// StreamDesc is an alias of grpc.StreamDesc.
	StreamDesc = grpc.StreamDesc

	// Streamer is an alias of grpc.Streamer.
	Streamer = grpc.Streamer

	// StreamHandler is an alias of grpc.StreamHandler.
	StreamHandler = grpc.StreamHandler

	// StreamServerInfo is an alias of grpc.StreamServerInfo.
	StreamServerInfo = grpc.StreamServerInfo

	// StreamServerInterceptor is an alias of grpc.StreamServerInterceptor.
	StreamServerInterceptor = grpc.StreamServerInterceptor

	// UnaryClientInterceptor is an alias of grpc.UnaryClientInterceptor.
	UnaryClientInterceptor = grpc.UnaryClientInterceptor

	// UnaryHandler is an alias of grpc.UnaryHandler.
	UnaryHandler = grpc.UnaryHandler

	// UnaryServerInfo is an alias of grpc.UnaryServerInfo.
	UnaryServerInfo = grpc.UnaryServerInfo

	// UnaryServerInterceptor is an alias of grpc.UnaryServerInterceptor.
	UnaryServerInterceptor = grpc.UnaryServerInterceptor
)

// ChainUnaryInterceptor is an alias of grpc.ChainUnaryInterceptor.
func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStreamInterceptor is an alias of grpc.ChainStreamInterceptor.
func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// Creds is an alias of grpc.Creds.
func Creds(c credentials.TransportCredentials) ServerOption {
	return grpc.Creds(c)
}

// NewClient is an alias of grpc.NewClient.
func NewClient(target string, opts ...DialOption) (*ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

// NewInsecureCredentials is an alias for insecure.NewCredentials.
func NewInsecureCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

// NewClient is an alias of credentials.NewTLS.
func NewTLS(c *tls.Config) credentials.TransportCredentials {
	return credentials.NewTLS(c)
}

// SetHeader is an alias of grpc.SetHeader.
func SetHeader(ctx context.Context, md metadata.MD) error {
	return grpc.SetHeader(ctx, md)
}

// TimeoutUnaryClientInterceptor is an alias of timeout.UnaryClientInterceptor.
func TimeoutUnaryClientInterceptor(d time.Duration) grpc.UnaryClientInterceptor {
	return timeout.UnaryClientInterceptor(d)
}

// UseCompressor is an alias of grpc.UseCompressor.
func UseCompressor(name string) CallOption {
	return grpc.UseCompressor(name)
}

// WithChainUnaryInterceptor is an alias of grpc.WithChainUnaryInterceptor.
func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// WithChainStreamInterceptor is an alias of grpc.WithChainStreamInterceptor.
func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

// WithContextDialer is an alias of grpc.WithContextDialer.
func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return grpc.WithContextDialer(f)
}

// WithDefaultCallOptions is an alias of grpc.WithDefaultCallOptions.
func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return grpc.WithDefaultCallOptions(cos...)
}

// WithUserAgent is an alias of grpc.WithUserAgent.
func WithUserAgent(s string) DialOption {
	return grpc.WithUserAgent(s)
}

// WithTransportCredentials is an alias of grpc.WithTransportCredentials.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return grpc.WithTransportCredentials(creds)
}

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
