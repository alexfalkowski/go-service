package token

import (
	"context"
	"fmt"
	"path"

	"github.com/alexfalkowski/go-service/security/header"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// ExtractToken from context.
func ExtractToken(ctx context.Context) (string, error) {
	md := meta.ExtractIncoming(ctx)

	values := md["authorization"]
	if len(values) == 0 {
		return "", header.ErrInvalidAuthorization
	}

	_, token, err := header.ParseAuthorization(values[0])

	return token, err
}

// VerifyToken from context.
func VerifyToken(ctx context.Context, verifier token.Verifier) (context.Context, error) {
	token, err := ExtractToken(ctx)
	if err != nil {
		return ctx, err
	}

	return verifier.Verify(ctx, []byte(token))
}

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(verifier token.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		ctx, err := VerifyToken(ctx, verifier)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(verifier token.Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		ctx := stream.Context()

		ctx, err := VerifyToken(ctx, verifier)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// NewPerRPCCredentials for token.
func NewPerRPCCredentials(generator token.Generator) credentials.PerRPCCredentials {
	return &tokenPerRPCCredentials{generator: generator}
}

type tokenPerRPCCredentials struct {
	generator token.Generator
}

func (p *tokenPerRPCCredentials) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	_, t, err := p.generator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, header.ErrInvalidAuthorization
	}

	return map[string]string{"authorization": fmt.Sprintf("%s %s", header.BearerAuthorization, t)}, nil
}

func (p *tokenPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
