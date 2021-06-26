package jwt

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/alexfalkowski/go-service/pkg/security/meta"
	grpcMeta "github.com/alexfalkowski/go-service/pkg/transport/grpc/meta"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	healthService = "grpc.health.v1.Health"
)

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(verifier jwt.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(ctx, req)
		}

		md := grpcMeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, jwt.ErrMissingToken.Error())
		}

		token, err := verifier.Verify(ctx, tkn(values[0]))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		ctx = meta.WithToken(ctx, token)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(verifier jwt.Verifier) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		md := grpcMeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, jwt.ErrMissingToken.Error())
		}

		token, err := verifier.Verify(ctx, tkn(values[0]))
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		wrapped := grpcMiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = meta.WithToken(ctx, token)

		return handler(srv, wrapped)
	}
}

// NewPerRPCCredentials for token.
func NewPerRPCCredentials(generator jwt.Generator) credentials.PerRPCCredentials {
	return &tokenPerRPCCredentials{generator: generator}
}

type tokenPerRPCCredentials struct {
	generator jwt.Generator
}

func (p *tokenPerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	t, err := p.generator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, jwt.ErrMissingToken
	}

	return map[string]string{"authorization": fmt.Sprintf("Bearer %s", t)}, nil
}

func (p *tokenPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}

func tkn(header string) []byte {
	s := strings.Split(header, " ")

	return []byte(s[1])
}
