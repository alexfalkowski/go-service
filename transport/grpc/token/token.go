package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
)

// NewAccessController returns an access controller when token auth is enabled.
//
// If cfg is disabled, it returns (nil, nil).
func NewAccessController(cfg *token.Config) (AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return access.NewController(cfg.Access)
}

// AccessController is an alias for access.Controller.
type AccessController access.Controller

// NewToken returns a token service when token auth is enabled.
//
// If cfg is disabled, it returns nil.
func NewToken(name env.Name, cfg *token.Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(name, cfg, fs, sig, ver, gen)}
}

// Token wraps token.Token for gRPC transport integration.
type Token struct {
	*token.Token
}

// NewVerifier returns a Verifier backed by token.
//
// If token is nil, it returns nil.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias for token.Verifier.
type Verifier token.Verifier

// UnaryServerInterceptor returns a gRPC unary server interceptor that verifies Authorization tokens.
//
// Requests with ignorable methods bypass verification.
// On verification failure, it returns Unauthenticated.
// On success, it stores the verified subject as the user id in the context and invokes the handler.
func UnaryServerInterceptor(id env.UserID, verifier Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(ctx, req)
		}

		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), info.FullMethod)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = meta.WithUserID(ctx, meta.Ignored(sub))
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that verifies Authorization tokens.
//
// Requests with ignorable methods bypass verification.
// On verification failure, it returns Unauthenticated.
// On success, it stores the verified subject as the user id in the stream context and invokes the handler.
func StreamServerInterceptor(id env.UserID, verifier Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), info.FullMethod)
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = meta.WithUserID(ctx, meta.Ignored(sub))
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// NewGenerator returns a Generator backed by token.
//
// If token is nil, it returns nil.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias for token.Generator.
type Generator token.Generator

// UnaryClientInterceptor returns a gRPC unary client interceptor that injects an Authorization token.
//
// It generates a token scoped to fullMethod and the provided user id and appends it to outgoing metadata
// under the "authorization" key as a Bearer token.
// On generation failure or an empty token, it returns Unauthenticated.
func UnaryClientInterceptor(id env.UserID, generator Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := generator.Generate(fullMethod, id.String())
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		if len(token) == 0 {
			return status.Error(codes.Unauthenticated, header.ErrInvalidAuthorization.Error())
		}

		auth := meta.Ignored(strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Append("authorization", auth.Value())

		ctx = meta.WithAuthorization(ctx, auth)
		ctx = meta.NewOutgoingContext(ctx, md)

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that injects an Authorization token.
//
// It generates a token scoped to fullMethod and the provided user id and appends it to outgoing metadata
// under the "authorization" key as a Bearer token.
// On generation failure or an empty token, it returns Unauthenticated.
func StreamClientInterceptor(id env.UserID, generator Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		token, err := generator.Generate(fullMethod, id.String())
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, header.ErrInvalidAuthorization.Error())
		}

		auth := meta.Ignored(strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Append("authorization", auth.Value())

		ctx = meta.WithAuthorization(ctx, auth)
		ctx = meta.NewOutgoingContext(ctx, md)

		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}
