package ratelimit

import (
	"context"

	"github.com/alexfalkowski/go-service/time/rate"
	"github.com/dgraph-io/ristretto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LimiterID to find the correct limiter. This way we can segment limiters by different ids.
type LimiterID func(ctx context.Context) string

// UnaryServerInterceptor for ratelimit.
func UnaryServerInterceptor(cfg *Config, cache *ristretto.Cache, limiterID LimiterID) grpc.UnaryServerInterceptor {
	params := rate.Params{Every: cfg.Every, Burst: cfg.Burst, TTL: cfg.TTL, Cache: cache}
	limiter := rate.New(params)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		l := limiter.Get(limiterID(ctx))

		if !l.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected as rate allowed is 1 req per %s, please retry later.", info.FullMethod, cfg.Every.String())
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for ratelimit.
func StreamServerInterceptor(cfg *Config, cache *ristretto.Cache, limiterID LimiterID) grpc.StreamServerInterceptor {
	params := rate.Params{Every: cfg.Every, Burst: cfg.Burst, TTL: cfg.TTL, Cache: cache}
	limiter := rate.New(params)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		l := limiter.Get(limiterID(ctx))

		if !l.Allow() {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected as rate allowed is 1 req per %s, please retry later.", info.FullMethod, cfg.Every.String())
		}

		return handler(srv, stream)
	}
}
