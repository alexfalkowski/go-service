package jwt

import (
	"context"
	"fmt"
	"path"
	"strings"

	sjwt "github.com/alexfalkowski/go-service/security/jwt"
	"github.com/alexfalkowski/go-service/security/meta"
	sstrings "github.com/alexfalkowski/go-service/strings"
	smeta "github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/golang-jwt/jwt/v4"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(verifier sjwt.Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if sstrings.IsHealth(service) {
			return handler(ctx, req)
		}

		md := smeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, sjwt.ErrMissingToken.Error())
		}

		token, err := verifier.Verify(ctx, tkn(values[0]))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		claims := token.Claims.(jwt.MapClaims)

		azp, ok := claims["azp"]
		if ok {
			ctx = meta.WithAuthorizedParty(ctx, azp.(string))
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(verifier sjwt.Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if sstrings.IsHealth(service) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		md := smeta.ExtractIncoming(ctx)

		values := md["authorization"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, sjwt.ErrMissingToken.Error())
		}

		token, err := verifier.Verify(ctx, tkn(values[0]))
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "could not verify token: %s", err.Error())
		}

		claims := token.Claims.(jwt.MapClaims)

		azp, ok := claims["azp"]
		if ok {
			ctx = meta.WithAuthorizedParty(ctx, azp.(string))
		}

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// NewPerRPCCredentials for token.
func NewPerRPCCredentials(generator sjwt.Generator) credentials.PerRPCCredentials {
	return &tokenPerRPCCredentials{generator: generator}
}

type tokenPerRPCCredentials struct {
	generator sjwt.Generator
}

func (p *tokenPerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	t, err := p.generator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, sjwt.ErrMissingToken
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
