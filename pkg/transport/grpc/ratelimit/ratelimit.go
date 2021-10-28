package ratelimit

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/time/rate"
	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserAgentUnaryServerInterceptor for ratelimit.
func UserAgentUnaryServerInterceptor(cfg *Config) grpc.UnaryServerInterceptor {
	limiter := rate.New(cfg.Every, cfg.Burst)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := limiter.Limiter(meta.UserAgent(ctx))

		if !l.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected as rate allowed is 1 req per %s, please retry later.", info.FullMethod, cfg.Every.String())
		}

		return handler(ctx, req)
	}
}

// UserAgentStreamServerInterceptor for ratelimit.
func UserAgentStreamServerInterceptor(cfg *Config) grpc.StreamServerInterceptor {
	limiter := rate.New(cfg.Every, cfg.Burst)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		l := limiter.Limiter(meta.UserAgent(ctx))

		if !l.Allow() {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected as rate allowed is 1 req per %s, please retry later.", info.FullMethod, cfg.Every.String())
		}

		return handler(srv, stream)
	}
}
