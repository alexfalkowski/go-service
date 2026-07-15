package grpc

import (
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/method"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor.
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/tap"
)

// CallOption is an alias of [grpc.CallOption].
//
// It represents a per-RPC option (for example compression, per-call
// credentials, etc.) passed to a client call.
type CallOption = grpc.CallOption

// ClientConn is an alias of [grpc.ClientConn].
//
// It represents a virtual connection to a gRPC endpoint and is used to create
// generated service clients.
type ClientConn = grpc.ClientConn

// ClientConnInterface is an alias of [grpc.ClientConnInterface].
//
// It represents the subset of client connection behavior required by generated
// service clients.
type ClientConnInterface = grpc.ClientConnInterface

// ClientStream is an alias of [grpc.ClientStream].
//
// It is the client-side stream interface used by streaming RPCs.
type ClientStream = grpc.ClientStream

// DialOption is an alias of [grpc.DialOption].
//
// It configures client connection creation (dialing, credentials, interceptors,
// stats handlers, etc.).
type DialOption = grpc.DialOption

// EmptyServerOption is an alias of [grpc.EmptyServerOption].
//
// It is used when building server options that conditionally return "no option".
type EmptyServerOption = grpc.EmptyServerOption

// MethodDesc is an alias of [grpc.MethodDesc].
//
// It describes a unary RPC method on a generated gRPC service.
type MethodDesc = grpc.MethodDesc

// UnaryInvoker is an alias of [grpc.UnaryInvoker].
//
// It is the function signature used by unary client interceptors to invoke the
// next interceptor/transport.
type UnaryInvoker = grpc.UnaryInvoker

// ServerOption is an alias of [grpc.ServerOption].
//
// It configures server construction (credentials, interceptors, keepalive,
// stats handlers, etc.).
type ServerOption = grpc.ServerOption

// Server is an alias of [grpc.Server].
//
// It is the gRPC server implementation that hosts registered services.
type Server = grpc.Server

// ServiceDesc is an alias of [grpc.ServiceDesc].
//
// It describes a generated gRPC service and its unary and stream methods.
type ServiceDesc = grpc.ServiceDesc

// ServerStream is an alias of [grpc.ServerStream].
//
// It is the server-side stream interface used by streaming RPCs.
type ServerStream = grpc.ServerStream

// ServiceRegistrar is an alias of [grpc.ServiceRegistrar].
//
// It is implemented by *[grpc.Server] and is accepted by generated Register*
// functions.
type ServiceRegistrar = grpc.ServiceRegistrar

// StreamClientInterceptor is an alias of [grpc.StreamClientInterceptor].
//
// It intercepts client-side streaming RPCs.
type StreamClientInterceptor = grpc.StreamClientInterceptor

// StreamDesc is an alias of [grpc.StreamDesc].
//
// It describes streaming RPC characteristics.
type StreamDesc = grpc.StreamDesc

// Streamer is an alias of [grpc.Streamer].
//
// It is the function signature used by stream client interceptors to create a
// client stream.
type Streamer = grpc.Streamer

// StreamHandler is an alias of [grpc.StreamHandler].
//
// It is the function signature used by stream server interceptors to handle a
// server stream.
type StreamHandler = grpc.StreamHandler

// StreamServerInfo is an alias of [grpc.StreamServerInfo].
//
// It provides information about a streaming RPC to a server interceptor.
type StreamServerInfo = grpc.StreamServerInfo

// StreamServerInterceptor is an alias of [grpc.StreamServerInterceptor].
//
// It intercepts server-side streaming RPCs.
type StreamServerInterceptor = grpc.StreamServerInterceptor

// UnaryClientInterceptor is an alias of [grpc.UnaryClientInterceptor].
//
// It intercepts client-side unary RPCs.
type UnaryClientInterceptor = grpc.UnaryClientInterceptor

// UnaryHandler is an alias of [grpc.UnaryHandler].
//
// It is the function signature used by unary server interceptors to invoke the
// handler.
type UnaryHandler = grpc.UnaryHandler

// UnaryServerInfo is an alias of [grpc.UnaryServerInfo].
//
// It provides information about a unary RPC to a server interceptor.
type UnaryServerInfo = grpc.UnaryServerInfo

// UnaryServerInterceptor is an alias of [grpc.UnaryServerInterceptor].
//
// It intercepts server-side unary RPCs.
type UnaryServerInterceptor = grpc.UnaryServerInterceptor

// TransportCredentials is an alias of [credentials.TransportCredentials].
type TransportCredentials = credentials.TransportCredentials

// TapInfo is an alias of [tap.Info].
//
// It carries per-RPC metadata (method name, wire byte count) available to a
// TapHandle before the request is read or dispatched to a handler.
type TapInfo = tap.Info

// TapHandle is an alias of [tap.ServerInHandle].
//
// It runs for every RPC as the connection accepts it, before message
// decoding or interceptors, and can reject the call by returning a non-nil
// error.
type TapHandle = tap.ServerInHandle

// ErrServerStopped is an alias for [grpc.ErrServerStopped].
var ErrServerStopped = grpc.ErrServerStopped

// MethodPolicy is an alias of [method.Policy].
//
// It stores gRPC method behavior used by server middleware.
type MethodPolicy = method.Policy

// NewMethodPolicy constructs a gRPC method policy.
func NewMethodPolicy() *MethodPolicy {
	return method.NewPolicy()
}

// StatusText returns the standard gRPC status text for the code.
func StatusText(code codes.Code) string {
	return codes.StatusText(code)
}

// ParseServiceMethod derives a logical service and method name from a gRPC full method.
//
// If name can be split as a slash-prefixed path, ParseServiceMethod returns the extracted service/method pair.
// Otherwise it returns "root" for both values.
func ParseServiceMethod(name string) (string, string) {
	if service, method, ok := url.SplitPath(name); ok {
		return service, method
	}

	return "root", "root"
}

// StatsHandler returns a ServerOption that installs h as the server stats handler.
//
// This is a thin wrapper around [grpc.StatsHandler]. A stats handler observes RPC
// lifecycle events (for example for telemetry) and is invoked by the gRPC
// server runtime.
func StatsHandler(h stats.Handler) ServerOption {
	return grpc.StatsHandler(h)
}

// InTapHandle returns a ServerOption that installs h as the server's tap handle.
//
// This is a thin wrapper around [grpc.InTapHandle]. A tap handle runs for every
// RPC as the connection accepts it, before message decoding or interceptors,
// and can reject the call early (for example for connection-level admission
// control or overload shedding) by returning a non-nil error.
func InTapHandle(h TapHandle) ServerOption {
	return grpc.InTapHandle(h)
}

// WithStatsHandler returns a DialOption that installs h as the client stats handler.
//
// This is a thin wrapper around [grpc.WithStatsHandler]. A stats handler observes
// outbound RPC lifecycle events (for example for telemetry) and is invoked by
// the gRPC client runtime.
func WithStatsHandler(h stats.Handler) DialOption {
	return grpc.WithStatsHandler(h)
}

// ChainUnaryInterceptor returns a ServerOption that chains unary server interceptors.
//
// It forwards to [grpc.ChainUnaryInterceptor]. Interceptors are executed in the
// order provided.
func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStreamInterceptor returns a ServerOption that chains stream server interceptors.
//
// It forwards to [grpc.ChainStreamInterceptor]. Interceptors are executed in the
// order provided.
func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// Creds returns a ServerOption that configures server-side transport credentials.
//
// It forwards to [grpc.Creds]. For TLS, use NewTLS to create
// [credentials.TransportCredentials] from a go-service [github.com/alexfalkowski/go-service/v2/crypto/tls.Config].
func Creds(c credentials.TransportCredentials) ServerOption {
	return grpc.Creds(c)
}

// MaxRecvMsgSize returns a ServerOption that caps inbound message size in bytes.
//
// This is a thin wrapper around [grpc.MaxRecvMsgSize].
func MaxRecvMsgSize(m int) ServerOption {
	return grpc.MaxRecvMsgSize(m)
}

// NewClient creates a new gRPC client channel for target using opts.
//
// This forwards to [grpc.NewClient]. No I/O is performed during construction; the
// returned ClientConn connects automatically when it is used for RPCs (or when
// Connect is called explicitly).
//
// The target format and supported schemes are defined by gRPC and any
// registered resolvers in your binary.
func NewClient(target string, opts ...DialOption) (*ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

// NewInsecureCredentials returns transport credentials that disable transport security.
//
// This is a thin wrapper around [insecure.NewCredentials]. Use this only for local
// development/testing or when transport security is provided out-of-band.
func NewInsecureCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

// NewTLS constructs TLS transport credentials from c.
//
// This is a thin wrapper around [credentials.NewTLS]. The provided
// [github.com/alexfalkowski/go-service/v2/crypto/tls.Config] is used by the gRPC transport for handshake and peer
// verification.
func NewTLS(c *tls.Config) credentials.TransportCredentials {
	return credentials.NewTLS(c)
}

// SetHeader sets the header metadata that will be sent back to the client.
//
// This forwards to [grpc.SetHeader]. It is typically called by a handler to add
// response headers.
func SetHeader(ctx context.Context, md meta.Map) error {
	return grpc.SetHeader(ctx, md)
}

// SetTrailer sets the trailer metadata that will be sent back to the client.
//
// This forwards to [grpc.SetTrailer]. It is typically called by a handler to add
// response trailers.
func SetTrailer(ctx context.Context, md meta.Map) error {
	return grpc.SetTrailer(ctx, md)
}

// Header returns a call option that captures response header metadata.
//
// This forwards to [grpc.Header]. The provided map is populated by the client
// call with any header metadata returned by the server.
func Header(md *meta.Map) CallOption {
	return grpc.Header(md)
}

// TimeoutUnaryClientInterceptor returns a unary client interceptor that applies a per-RPC timeout.
//
// This forwards to [timeout.UnaryClientInterceptor] from go-grpc-middleware.
// The interceptor derives every outgoing context with a timeout of d. If the
// parent already has an earlier deadline, that earlier deadline remains in
// effect; the derived timeout cannot extend it.
func TimeoutUnaryClientInterceptor(d time.Duration) grpc.UnaryClientInterceptor {
	return timeout.UnaryClientInterceptor(d.Duration())
}

// TimeoutUnaryServerInterceptor returns a unary server interceptor that applies a per-RPC timeout.
//
// The interceptor wraps the incoming context with d before invoking the next handler.
// Existing earlier deadlines are preserved because the derived context cannot extend
// the parent context's deadline.
func TimeoutUnaryServerInterceptor(d time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		timedCtx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		return handler(timedCtx, req)
	}
}

// UseCompressor returns a CallOption that requests message compression by name.
//
// This forwards to [grpc.UseCompressor]. The compressor must be registered with
// gRPC (for example gzip is installed by importing
// google.golang.org/grpc/encoding/gzip).
func UseCompressor(name string) CallOption {
	return grpc.UseCompressor(name)
}

// WithChainUnaryInterceptor returns a DialOption that chains unary client interceptors.
//
// This forwards to [grpc.WithChainUnaryInterceptor]. Interceptors are executed in
// the order provided.
func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// WithChainStreamInterceptor returns a DialOption that chains stream client interceptors.
//
// This forwards to [grpc.WithChainStreamInterceptor]. Interceptors are executed in
// the order provided.
func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

// WithContextDialer returns a DialOption that sets the dial function used to create connections.
//
// This forwards to [grpc.WithContextDialer]. The dialer is invoked with the
// outgoing context and the resolved address.
func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption {
	return grpc.WithContextDialer(f)
}

// WithDefaultCallOptions returns a DialOption that configures default CallOptions.
//
// This forwards to [grpc.WithDefaultCallOptions]. The provided options are applied
// to all RPCs created from the resulting connection unless overridden per-call.
func WithDefaultCallOptions(cos ...CallOption) DialOption {
	return grpc.WithDefaultCallOptions(cos...)
}

// WithUserAgent returns a DialOption that sets the user-agent string for the connection.
//
// This forwards to [grpc.WithUserAgent].
func WithUserAgent(s string) DialOption {
	return grpc.WithUserAgent(s)
}

// WithTransportCredentials returns a DialOption that configures client-side transport credentials.
//
// This forwards to [grpc.WithTransportCredentials]. For TLS, use NewTLS to
// create [credentials.TransportCredentials] from a go-service
// [github.com/alexfalkowski/go-service/v2/crypto/tls.Config].
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return grpc.WithTransportCredentials(creds)
}

// WithKeepaliveParams configures client keepalive ping interval and timeout.
//
// This forwards to [grpc.WithKeepaliveParams], setting PermitWithoutStream to true
// so that the client may send pings even when there are no active RPC streams.
func WithKeepaliveParams(ping, timeout time.Duration) DialOption {
	return grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                ping.Duration(),
		Timeout:             timeout.Duration(),
		PermitWithoutStream: true,
	})
}

// NewServer constructs a *[grpc.Server] with standard keepalive configuration and reflection enabled.
//
// Keepalive enforcement, connection establishment timeout, and server parameters are
// populated from options using the following duration keys:
//
//   - keepalive_enforcement_policy_ping_min_time (falls back to timeout)
//   - keepalive_max_connection_idle (falls back to timeout)
//   - keepalive_max_connection_age (falls back to the gRPC default)
//   - keepalive_max_connection_age_grace (falls back to the gRPC default)
//   - keepalive_ping_time (falls back to timeout)
//   - connection_timeout (falls back to timeout)
//
// In addition, the provided timeout is used as the keepalive ping Timeout.
//
// Additional low-level server tuning may be provided through options using:
//
//   - max_concurrent_streams
//   - max_header_list_size
//   - initial_window_size
//   - initial_conn_window_size
//   - max_send_msg_size
//
// This function always registers server reflection via [reflection.Register] so
// tools such as grpcurl can discover services. Deployments that should not expose
// reflection publicly should restrict access with bind addresses, TLS/client
// authentication, ingress policy, firewall rules, or service-mesh authorization.
//
// Any additional opts are appended after the keepalive options and may further
// customize server behavior (for example interceptors, credentials, or stats handlers).
func NewServer(options options.Map, timeout time.Duration, opts ...ServerOption) *Server {
	serverOptions := make([]ServerOption, 0, 8+len(opts))
	serverOptions = append(serverOptions, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             options.NonNegativeDuration("keepalive_enforcement_policy_ping_min_time", timeout).Duration(),
		PermitWithoutStream: true,
	}))
	serverOptions = append(serverOptions, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     options.NonNegativeDuration("keepalive_max_connection_idle", timeout).Duration(),
		MaxConnectionAge:      options.NonNegativeDuration("keepalive_max_connection_age", 0).Duration(),
		MaxConnectionAgeGrace: options.NonNegativeDuration("keepalive_max_connection_age_grace", 0).Duration(),
		Time:                  options.NonNegativeDuration("keepalive_ping_time", timeout).Duration(),
		Timeout:               timeout.Duration(),
	}))
	serverOptions = append(serverOptions, grpc.ConnectionTimeout(options.NonNegativeDuration("connection_timeout", timeout).Duration()))
	if _, ok := options["max_concurrent_streams"]; ok {
		serverOptions = append(serverOptions, grpc.MaxConcurrentStreams(options.Uint32("max_concurrent_streams", 0)))
	}

	if _, ok := options["max_header_list_size"]; ok {
		serverOptions = append(serverOptions, grpc.MaxHeaderListSize(options.Uint32Size("max_header_list_size", 0)))
	}

	if _, ok := options["initial_window_size"]; ok {
		serverOptions = append(serverOptions, grpc.InitialWindowSize(options.Int32Size("initial_window_size", 0)))
	}

	if _, ok := options["initial_conn_window_size"]; ok {
		serverOptions = append(serverOptions, grpc.InitialConnWindowSize(options.Int32Size("initial_conn_window_size", 0)))
	}

	if _, ok := options["max_send_msg_size"]; ok {
		serverOptions = append(serverOptions, grpc.MaxSendMsgSize(options.IntSize("max_send_msg_size", 0)))
	}
	serverOptions = append(serverOptions, opts...)

	server := grpc.NewServer(serverOptions...)
	// security: reflection is intentionally always enabled for go-service
	// servers so internal tooling can discover registered services. Restrict
	// public exposure at the network/auth boundary when needed.
	reflection.Register(server)

	return server
}
