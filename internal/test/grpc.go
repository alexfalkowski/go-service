package test

import (
	"maps"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/test/bufconn"
)

// BufconnOption configures an in-memory gRPC transport test connection.
type BufconnOption func(*bufconnOptions)

type bufconnOptions struct {
	generator            token.Generator
	verifier             token.Verifier
	serverLimiter        *limiter.Server
	clientLimiter        *limiter.Client
	logger               *logger.Logger
	serverOptions        options.Map
	clientBreakerOptions []breaker.Option
	maxReceiveSize       bytes.Size
	clientBreaker        bool
	compression          bool
}

// WithBufconnGenerator configures client-side token generation for the in-memory gRPC connection.
func WithBufconnGenerator(generator token.Generator) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.generator = generator
	}
}

// WithBufconnLogger configures client and server logging for the in-memory gRPC connection.
func WithBufconnLogger(logger *logger.Logger) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.logger = logger
	}
}

// WithBufconnVerifier configures server-side token verification for the in-memory gRPC connection.
func WithBufconnVerifier(verifier token.Verifier) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.verifier = verifier
	}
}

// WithBufconnServerLimiter configures server-side rate limiting for the in-memory gRPC connection.
func WithBufconnServerLimiter(limiter *limiter.Server) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.serverLimiter = limiter
	}
}

// WithBufconnClientLimiter configures client-side rate limiting for the in-memory gRPC connection.
func WithBufconnClientLimiter(limiter *limiter.Client) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.clientLimiter = limiter
	}
}

// WithBufconnBreaker configures client-side circuit breaking for the in-memory gRPC connection.
func WithBufconnBreaker(opts ...breaker.Option) BufconnOption {
	return func(options *bufconnOptions) {
		options.clientBreaker = true
		options.clientBreakerOptions = opts
	}
}

// WithBufconnCompression enables client-side compression for the in-memory gRPC connection.
func WithBufconnCompression() BufconnOption {
	return func(opts *bufconnOptions) {
		opts.compression = true
	}
}

// WithBufconnMaxReceiveSize configures the server receive size for the in-memory gRPC connection.
func WithBufconnMaxReceiveSize(size bytes.Size) BufconnOption {
	return func(opts *bufconnOptions) {
		opts.maxReceiveSize = size
	}
}

// WithBufconnServerOptions configures low-level server options for the in-memory gRPC connection.
func WithBufconnServerOptions(opts options.Map) BufconnOption {
	return func(options *bufconnOptions) {
		options.serverOptions = opts
	}
}

// NewBufconnGRPCConn returns a gRPC client connection backed by bufconn.
func NewBufconnGRPCConn(tb testing.TB, opts ...BufconnOption) *grpc.ClientConn {
	tb.Helper()

	os := &bufconnOptions{serverOptions: options.Map{}}
	for _, opt := range opts {
		opt(os)
	}
	if os.logger == nil {
		lc := fxtest.NewLifecycle(tb)
		logger, err := NewLogger(lc, NewTextLoggerConfig())
		require.NoError(tb, err)
		os.logger = logger
	}

	listener := bufconn.Listen(1024 * 1024)
	server := newBufconnGRPCServer(os)
	v1.RegisterGreeterServiceServer(server, NewService())

	go func() {
		_ = server.Serve(listener)
	}()

	tb.Cleanup(func() {
		server.GracefulStop()
		_ = listener.Close()
	})

	clientOpts := []transportgrpc.ClientOption{
		transportgrpc.WithClientUnaryInterceptors(),
		transportgrpc.WithClientStreamInterceptors(),
		transportgrpc.WithClientUserAgent(UserAgent),
		transportgrpc.WithClientID(uuid.NewGenerator()),
		transportgrpc.WithClientDialOption(grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		})),
		transportgrpc.WithClientLogger(os.logger),
		transportgrpc.WithClientLimiter(os.clientLimiter),
	}

	if os.generator != nil {
		clientOpts = append(clientOpts, transportgrpc.WithClientTokenGenerator(UserID, os.generator))
	}
	if os.clientBreaker {
		clientOpts = append(clientOpts, transportgrpc.WithClientBreaker(os.clientBreakerOptions...))
	}
	if os.compression {
		clientOpts = append(clientOpts, transportgrpc.WithClientCompression())
	}

	conn, err := transportgrpc.NewClient("passthrough:///bufnet", clientOpts...)
	require.NoError(tb, err)

	tb.Cleanup(func() {
		_ = conn.Close()
	})

	return conn
}

func newBufconnGRPCServer(opts *bufconnOptions) *grpc.Server {
	id := uuid.NewGenerator()

	unary := []grpc.UnaryServerInterceptor{
		meta.UnaryServerInterceptor(UserAgent, Version, id),
		logger.UnaryServerInterceptor(opts.logger),
	}
	if opts.verifier != nil {
		unary = append(unary, token.UnaryServerInterceptor(UserID, opts.verifier))
	}
	if opts.serverLimiter != nil {
		unary = append(unary, limiter.UnaryServerInterceptor(opts.serverLimiter))
	}

	stream := []grpc.StreamServerInterceptor{
		meta.StreamServerInterceptor(UserAgent, Version, id),
		logger.StreamServerInterceptor(opts.logger),
	}
	if opts.verifier != nil {
		stream = append(stream, token.StreamServerInterceptor(UserID, opts.verifier))
	}

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unary...),
		grpc.ChainStreamInterceptor(stream...),
	}
	if opts.maxReceiveSize > 0 {
		serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(int(opts.maxReceiveSize.Bytes())))
	}

	os := options.Map{}
	maps.Copy(os, opts.serverOptions)

	return grpc.NewServer(os, DefaultTimeout, serverOptions...)
}
