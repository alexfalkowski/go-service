package token

import (
	"context"
	"fmt"
	"path"
	"strings"

	grpcMeta "github.com/alexfalkowski/go-service/pkg/grpc/meta"
	"github.com/alexfalkowski/go-service/pkg/security/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	healthService = "grpc.health.v1.Health"
)

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(verifier token.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(ctx, req)
		}

		md := grpcMeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		if err := verifier.Verify(tkn(values[0])); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(verifier token.Verifier) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		md := grpcMeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		if err := verifier.Verify(tkn(values[0])); err != nil {
			return status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		return handler(srv, stream)
	}
}

// NewPerRPCCredentials for token.
func NewPerRPCCredentials(generator token.Generator) credentials.PerRPCCredentials {
	return &tokenPerRPCCredentials{generator: generator}
}

type tokenPerRPCCredentials struct {
	generator token.Generator
}

func (p *tokenPerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := p.generator.Generate()
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return map[string]string{}, nil
	}

	return map[string]string{"authorization": fmt.Sprintf("Bearer %s", token)}, nil
}

func (p *tokenPerRPCCredentials) RequireTransportSecurity() bool {
	return false
}

func tkn(header string) []byte {
	s := strings.Split(header, " ")

	return []byte(s[1])
}
