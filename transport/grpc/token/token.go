package token

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	ts "github.com/alexfalkowski/go-service/v2/transport/strings"
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
		if ts.IsObservable(service) {
			return handler(ctx, req)
		}

		token := meta.Authorization(ctx).Value()

		ctx, err := verifier.Verify(ctx, strings.Bytes(token))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(verifier token.Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if ts.IsObservable(service) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		token := meta.Authorization(ctx).Value()

		ctx, err := verifier.Verify(ctx, strings.Bytes(token))
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

	meta := map[string]string{
		"authorization": strings.Join(" ", header.BearerAuthorization, bytes.String(token)),
	}

	return meta, nil
}

func (p *tokenPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}
