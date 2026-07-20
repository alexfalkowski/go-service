package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/method"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
)

// NewToken constructs a token service when token auth is enabled.
//
// The returned service is responsible for generating and verifying tokens according to cfg (for example,
// JWT/PASETO/SSH token kinds as configured by the underlying token package).
//
// If cfg is disabled, it returns nil so callers can treat token auth as not configured.
func NewToken(cfg *token.Config, fs *os.FS, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(cfg, fs, gen)}
}

// Token wraps *[github.com/alexfalkowski/go-service/v2/token.Token] for gRPC transport integration.
//
// It exists so transport-level wiring can keep a distinct type for gRPC token functionality while still
// delegating generation and verification to the underlying token implementation.
type Token struct {
	*token.Token
}

// NewVerifier returns a [Verifier] backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a verifier only when token auth
// is enabled/configured, and to leave verification interceptors disabled otherwise.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias for [token.Verifier].
//
// Verifiers validate Authorization tokens and typically return a "subject" string (the authenticated
// principal) on success.
type Verifier token.Verifier

// UnaryServerInterceptor returns a gRPC unary server interceptor that verifies Authorization tokens.
//
// Operation and unauthenticated methods bypass verification.
//
// The interceptor expects an Authorization value to have been extracted into the context by the metadata
// interceptor ([github.com/alexfalkowski/go-service/v2/net/grpc/meta.UnaryServerInterceptor]). It verifies the token using verifier, scoping
// verification to the RPC `FullMethod`.
//
// Behavior:
//   - If verification fails, it returns [codes.Unauthenticated].
//   - If verification succeeds, it stores the verified subject as the user id in the context and invokes
//     the handler.
//
// Callers should only install this interceptor when verifier is non-nil.
func UnaryServerInterceptor(policy *method.Policy, verifier Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if bypassAuth(policy, info.FullMethod) {
			return handler(ctx, req)
		}

		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), info.FullMethod)
		if err != nil {
			return nil, status.SafeError(codes.Unauthenticated, err)
		}

		ctx = meta.WithAttributes(ctx, meta.WithUserID(meta.Ignored(sub)))
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that verifies Authorization tokens.
//
// Operation and unauthenticated methods bypass verification.
//
// The interceptor expects an Authorization value to have been extracted into the stream context by the
// metadata interceptor ([github.com/alexfalkowski/go-service/v2/net/grpc/meta.StreamServerInterceptor]). It verifies the token using verifier,
// scoping verification to the RPC `FullMethod`.
//
// Behavior:
//   - If verification fails, it returns [codes.Unauthenticated].
//   - If verification succeeds, it injects the verified subject as the user id into the stream context and
//     invokes the handler using a wrapped stream (`go-grpc-middleware` wrapper) that carries the new context.
//
// Callers should only install this interceptor when verifier is non-nil.
func StreamServerInterceptor(policy *method.Policy, verifier Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if bypassAuth(policy, info.FullMethod) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), info.FullMethod)
		if err != nil {
			return status.SafeError(codes.Unauthenticated, err)
		}

		ctx = meta.WithAttributes(ctx, meta.WithUserID(meta.Ignored(sub)))
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// UnaryAccessServerInterceptor returns a gRPC unary server interceptor that enforces access policy.
//
// Operation and unauthenticated methods bypass access control. For application RPCs, a missing verified user id is treated
// as unauthenticated, a policy denial returns [codes.PermissionDenied], and policy evaluation errors return
// [codes.Internal].
func UnaryAccessServerInterceptor(policy *method.Policy, controller access.Controller) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if bypassAuth(policy, info.FullMethod) {
			return handler(ctx, req)
		}

		if meta.UserID(ctx).IsEmpty() {
			return nil, status.SafeError(codes.Unauthenticated, header.ErrInvalidAuthorization)
		}

		ok, err := controller.HasAccess(ctx)
		if err != nil {
			return nil, status.SafeError(codes.Internal, err)
		}
		if !ok {
			return nil, status.SafeError(codes.PermissionDenied, access.ErrAccessDenied)
		}

		return handler(ctx, req)
	}
}

// StreamAccessServerInterceptor returns a gRPC stream server interceptor that enforces access policy.
//
// Operation and unauthenticated methods bypass access control. For application RPCs, a missing verified user id is treated
// as unauthenticated, a policy denial returns [codes.PermissionDenied], and policy evaluation errors return
// [codes.Internal].
func StreamAccessServerInterceptor(policy *method.Policy, controller access.Controller) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if bypassAuth(policy, info.FullMethod) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		if meta.UserID(ctx).IsEmpty() {
			return status.SafeError(codes.Unauthenticated, header.ErrInvalidAuthorization)
		}

		ok, err := controller.HasAccess(ctx)
		if err != nil {
			return status.SafeError(codes.Internal, err)
		}
		if !ok {
			return status.SafeError(codes.PermissionDenied, access.ErrAccessDenied)
		}

		return handler(srv, stream)
	}
}

func bypassAuth(policy *method.Policy, fullMethod string) bool {
	return policy.IsOperation(fullMethod) || policy.IsUnauthenticated(fullMethod)
}

// NewGenerator returns a [Generator] backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a token generator only when token
// auth is enabled/configured, and to leave client token-injection interceptors disabled otherwise.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias for [token.Generator].
//
// Generators create Authorization tokens for outbound RPCs, typically scoped to the RPC full method name
// and a caller identity (user id).
type Generator token.Generator

// UnaryClientInterceptor returns a gRPC unary client interceptor that injects an Authorization token.
//
// For each outbound unary RPC, it generates a token scoped to `fullMethod` and the provided user id and
// stores it in outgoing metadata under the "authorization" key using the `Bearer` scheme.
//
// The interceptor also stores the Authorization value in the context via [meta.WithAttributes], which can
// be useful for downstream instrumentation. Any existing outgoing "authorization" metadata is replaced so
// stale values do not take precedence over the newly generated token.
//
// Failure behavior:
//   - If token generation fails, it returns [codes.Unauthenticated].
//   - If token generation returns an empty token, it returns [codes.Unauthenticated] with
//     [header.ErrInvalidAuthorization].
//
// Callers should only install this interceptor when generator is non-nil.
func UnaryClientInterceptor(id env.UserID, generator Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := generator.Generate(fullMethod, id.String())
		if err != nil {
			return status.SafeError(codes.Unauthenticated, err)
		}

		if len(token) == 0 {
			return status.SafeError(codes.Unauthenticated, header.ErrInvalidAuthorization)
		}

		auth := meta.Ignored(strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Set("authorization", auth.Value())

		ctx = meta.WithAttributes(ctx, meta.WithAuthorization(auth))
		ctx = meta.NewOutgoingContext(ctx, md)

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that injects an Authorization token.
//
// For each outbound streaming RPC, it generates a token scoped to `fullMethod` and the provided user id and
// stores it in outgoing metadata under the "authorization" key using the `Bearer` scheme.
//
// The interceptor also stores the Authorization value in the context via [meta.WithAttributes]. Any
// existing outgoing "authorization" metadata is replaced so stale values do not take precedence over the
// newly generated token.
//
// Failure behavior:
//   - If token generation fails, it returns [codes.Unauthenticated].
//   - If token generation returns an empty token, it returns [codes.Unauthenticated] with
//     [header.ErrInvalidAuthorization].
//
// Callers should only install this interceptor when generator is non-nil.
func StreamClientInterceptor(id env.UserID, generator Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		token, err := generator.Generate(fullMethod, id.String())
		if err != nil {
			return nil, status.SafeError(codes.Unauthenticated, err)
		}

		if len(token) == 0 {
			return nil, status.SafeError(codes.Unauthenticated, header.ErrInvalidAuthorization)
		}

		auth := meta.Ignored(strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Set("authorization", auth.Value())

		ctx = meta.WithAttributes(ctx, meta.WithAuthorization(auth))
		ctx = meta.NewOutgoingContext(ctx, md)

		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}
