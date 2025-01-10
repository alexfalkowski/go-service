package token

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport/header"
	"github.com/alexfalkowski/go-service/transport/strings"
	tt "github.com/alexfalkowski/go-service/transport/token"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(verifier token.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(ctx, req)
		}

		ctx, err := tt.Verify(ctx, verifier)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
//
//nolint:fatcontext
func StreamServerInterceptor(verifier token.Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(srv, stream)
		}

		ctx := stream.Context()

		ctx, err := tt.Verify(ctx, verifier)
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
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
	_, token, err := p.generator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return nil, header.ErrInvalidAuthorization
	}

	return map[string]string{"authorization": header.BearerAuthorization + " " + string(token)}, nil
}

func (p *tokenPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
