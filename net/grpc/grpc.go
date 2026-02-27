package grpc

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/config/options"
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
	"google.golang.org/grpc/stats"
)

type (
	// CallOption is an alias of grpc.CallOption.
	//
	// It represents a per-RPC option (for example compression, per-call
	// credentials, etc.) passed to a client call.
	CallOption = grpc.CallOption

	// ClientConn is an alias of grpc.ClientConn.
	//
	// It represents a virtual connection to a gRPC endpoint and is used to create
	// generated service clients.
	ClientConn = grpc.ClientConn

	// ClientStream is an alias of grpc.ClientStream.
	//
	// It is the client-side stream interface used by streaming RPCs.
	ClientStream = grpc.ClientStream

	// DialOption is an alias of grpc.DialOption.
	//
	// It configures client connection creation (dialing, credentials, interceptors,
	// stats handlers, etc.).
	DialOption = grpc.DialOption

	// EmptyServerOption is an alias of grpc.EmptyServerOption.
	//
	// It is used when building server options that conditionally return "no option".
	EmptyServerOption = grpc.EmptyServerOption

	// UnaryInvoker is an alias of grpc.UnaryInvoker.
	//
	// It is the function signature used by unary client interceptors to invoke the
	// next interceptor/transport.
	UnaryInvoker = grpc.UnaryInvoker

	// ServerOption is an alias of grpc.ServerOption.
	//
	// It configures server construction (credentials, interceptors, keepalive,
	// stats handlers, etc.).
	ServerOption = grpc.ServerOption

	// Server is an alias of grpc.Server.
	//
	// It is the gRPC server implementation that hosts registered services.
	Server = grpc.Server

	// ServerStream is an alias of grpc.ServerStream.
	//
	// It is the server-side stream interface used by streaming RPCs.
	ServerStream = grpc.ServerStream

	// ServiceRegistrar is an alias of grpc.ServiceRegistrar.
	//
	// It is implemented by *grpc.Server and is accepted by generated Register*
	// functions.
	ServiceRegistrar = grpc.ServiceRegistrar

	// StreamClientInterceptor is an alias of grpc.StreamClientInterceptor.
	//
	// It intercepts client-side streaming RPCs.
	StreamClientInterceptor = grpc.StreamClientInterceptor

	// StreamDesc is an alias of grpc.StreamDesc.
	//
	// It describes streaming RPC characteristics.
	StreamDesc = grpc.StreamDesc

	// Streamer is an alias of grpc.Streamer.
	//
	// It is the function signature used by stream client interceptors to create a
	// client stream.
	Streamer = grpc.Streamer

	// StreamHandler is an alias of grpc.StreamHandler.
	//
	// It is the function signature used by stream server interceptors to handle a
	// server stream.
	StreamHandler = grpc.StreamHandler

	// StreamServerInfo is an alias of grpc.StreamServerInfo.
	//
	// It provides information about a streaming RPC to a server interceptor.
	StreamServerInfo = grpc.StreamServerInfo

	// StreamServerInterceptor is an alias of grpc.StreamServerInterceptor.
	//
	// It intercepts server-side streaming RPCs.
	StreamServerInterceptor = grpc.StreamServerInterceptor

	// UnaryClientInterceptor is an alias of grpc.UnaryClientInterceptor.
	//
	// It intercepts client-side unary RPCs.
	UnaryClientInterceptor = grpc.UnaryClientInterceptor

	// UnaryHandler is an alias of grpc.UnaryHandler.
	//
	// It is the function signature used by unary server interceptors to invoke the
	// handler.
	UnaryHandler = grpc.UnaryHandler

	// UnaryServerInfo is an alias of grpc.UnaryServerInfo.
	//
	// It provides information about a unary RPC to a server interceptor.
	UnaryServerInfo = grpc.UnaryServerInfo

	// UnaryServerInterceptor is an alias of grpc.UnaryServerInterceptor.
	//
	// It intercepts server-side unary RPCs.
	UnaryServerInterceptor = grpc.UnaryServerInterceptor
)

// StatsHandler returns a ServerOption that installs h as the server stats handler.
//
// This is a thin wrapper around grpc.StatsHandler. A stats handler observes RPC
// lifecycle events (for example for telemetry) and is invoked by the gRPC
// server runtime.
func StatsHandler(h stats.Handler) ServerOption {
	return grpc.StatsHandler(h)
}

// WithStatsHandler returns a DialOption that installs h as the client stats handler.
//
// This is a thin wrapper around grpc.WithStatsHandler. A stats handler observes
// outbound RPC lifecycle events (for example for telemetry) and is invoked by
// the gRPC client runtime.
func WithStatsHandler(h stats.Handler) DialOption {
	return grpc.WithStatsHandler(h)
}

// ChainUnaryInterceptor returns a ServerOption that chains unary server interceptors.
//
// It forwards to grpc.ChainUnaryInterceptor. Interceptors are executed in the
// order provided.
func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStreamInterceptor returns a ServerOption that chains stream server interceptors.
//
// It forwards to grpc.ChainStreamInterceptor. Interceptors are executed in the
// order provided.
func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// Creds returns a ServerOption that configures server-side transport credentials.
//
// It forwards to grpc.Creds. For TLS, use NewTLS to create credentials.TransportCredentials
// from a *tls.Config.
func Creds(c credentials.TransportCredentials) ServerOption {
	return grpc.Creds(c)
}

// NewClient creates a gRPC client connection to target using opts.
//
// This forwards to grpc.NewClient. The target format and supported schemes are
// defined by gRPC and any registered resolvers in your binary.
func NewClient(target string, opts ...DialOption) (*ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

// NewInsecureCredentials returns transport credentials that disable transport security.
//
// This is a thin wrapper around insecure.NewCredentials. Use this only for local
// development/testing or when transport security is provided out-of-band.
func NewInsecureCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

// NewTLS constructs TLS transport credentials from c.
//
// This is a thin wrapper around credentials.NewTLS. The provided tls.Config is
// used by the gRPC transport for handshake and peer verification.
func NewTLS(c *tls.Config) credentials.TransportCredentials {
	return credentials.NewTLS(c)
}

// SetHeader sets the header metadata that will be sent back to the client.
//
// This forwards to grpc.SetHeader. It is typically called by a handler to add
// response headers.
func SetHeader(ctx context.Context, md metadata.MD) error {
	return grpc.SetHeader(ctx, md)
}

// TimeoutUnaryClientInterceptor returns a unary client interceptor that applies a per-RPC timeout.
//
// This forwards to timeout.UnaryClientInterceptor from go-grpc-middleware.
// The interceptor typically wraps the outgoing context with a deadline of d.
func TimeoutUnaryClientInterceptor(d time.Duration) grpc.UnaryClientInterceptor {
	return timeout.UnaryClientInterceptor(d)
}

// UseCompressor returns a CallOption that requests message compression by name.
//
// This forwards to grpc.UseCompressor. The compressor must be registered with
// gRPC (for example gzip is installed by importing
// google.golang.org/grpc/encoding/gzip).
func UseCompressor(name string) CallOption {
	return grpc.UseCompressor(name)
}

// WithChainUnaryInterceptor returns a DialOption that chains unary client interceptors.
//
// This forwards to grpc.WithChainUnaryInterceptor. Interceptors are executed in
// the order provided.
func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// WithChainStreamInterceptor returns a DialOption that chains stream client interceptors.
//
// This forwards to grpc.WithChainStreamInterceptor. Interceptors are executed in
// the order provided.
func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

// WithContextDialer returns a DialOption that sets the dial function used to create connections.
//
// This forwards to grpc.WithContextDialer. The dialer is invoked with the
// outgoing context and the resolved address.
func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return grpc.WithContextDialer(f)
}

// WithDefaultCallOptions returns a DialOption that configures default CallOptions.
//
// This forwards to grpc.WithDefaultCallOptions. The provided options are applied
// to all RPCs created from the resulting connection unless overridden per-call.
func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return grpc.WithDefaultCallOptions(cos...)
}

// WithUserAgent returns a DialOption that sets the user-agent string for the connection.
//
// This forwards to grpc.WithUserAgent.
func WithUserAgent(s string) DialOption {
	return grpc.WithUserAgent(s)
}

// WithTransportCredentials returns a DialOption that configures client-side transport credentials.
//
// This forwards to grpc.WithTransportCredentials. For TLS, use NewTLS to create
// credentials.TransportCredentials from a *tls.Config.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return grpc.WithTransportCredentials(creds)
}

// WithKeepaliveParams configures client keepalive ping interval and timeout.
//
// This forwards to grpc.WithKeepaliveParams, setting PermitWithoutStream to true
// so that the client may send pings even when there are no active RPC streams.
func WithKeepaliveParams(ping, timeout time.Duration) DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                ping,
		Timeout:             timeout,
		PermitWithoutStream: true,
	})
}

// NewServer constructs a *grpc.Server with standard keepalive configuration and reflection enabled.
//
// Keepalive enforcement and server parameters are populated from options using the
// following duration keys (falling back to timeout when a key is not present):
//
//   - keepalive_enforcement_policy_ping_min_time
//   - keepalive_max_connection_idle
//   - keepalive_max_connection_age
//   - keepalive_max_connection_age_grace
//   - keepalive_ping_time
//
// In addition, the provided timeout is used as the keepalive ping Timeout.
//
// This function always registers server reflection via reflection.Register so
// tools such as grpcurl can discover services when reflection is permitted.
//
// Any additional opts are appended after the keepalive options and may further
// customize server behavior (for example interceptors, credentials, or stats handlers).
func NewServer(options options.Map, timeout time.Duration, opts ...ServerOption) *Server {
	os := make([]ServerOption, 0, 2+len(opts))
	os = append(os, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             options.Duration("keepalive_enforcement_policy_ping_min_time", timeout),
		PermitWithoutStream: true,
	}))
	os = append(os, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     options.Duration("keepalive_max_connection_idle", timeout),
		MaxConnectionAge:      options.Duration("keepalive_max_connection_age", timeout),
		MaxConnectionAgeGrace: options.Duration("keepalive_max_connection_age_grace", timeout),
		Time:                  options.Duration("keepalive_ping_time", timeout),
		Timeout:               timeout,
	}))
	os = append(os, opts...)

	server := grpc.NewServer(os...)
	reflection.Register(server)

	return server
}
