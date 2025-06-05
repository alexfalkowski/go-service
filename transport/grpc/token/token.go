package token

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	// Generator is an alias token.Generator.
	Generator = token.Generator

	// Verifier is an alias token.Verifier.
	Verifier = token.Verifier
)

// UnaryServerInterceptor for token.
func UnaryServerInterceptor(id env.UserID, verifier Verifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), p)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = meta.WithUserID(ctx, meta.Ignored(sub))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for token.
func StreamServerInterceptor(id env.UserID, verifier Verifier) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		auth := meta.Authorization(ctx).Value()

		sub, err := verifier.Verify(strings.Bytes(auth), p)
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = meta.WithUserID(ctx, meta.Ignored(sub))

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

// UnaryClientInterceptor for token.
func UnaryClientInterceptor(id env.UserID, generator Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := fullMethod[1:]

		token, err := generator.Generate(p, id.String())
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		if len(token) == 0 {
			return status.Error(codes.Unauthenticated, header.ErrInvalidAuthorization.Error())
		}

		auth := meta.Ignored(strings.Join(" ", header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Append("authorization", auth.Value())

		ctx = meta.WithAuthorization(ctx, auth)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor for token.
func StreamClientInterceptor(id env.UserID, generator Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := fullMethod[1:]

		token, err := generator.Generate(p, id.String())
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, header.ErrInvalidAuthorization.Error())
		}

		auth := meta.Ignored(strings.Join(" ", header.BearerAuthorization, bytes.String(token)))

		md := meta.ExtractOutgoing(ctx)
		md.Append("authorization", auth.Value())

		ctx = meta.WithAuthorization(ctx, auth)
		ctx = metadata.NewOutgoingContext(ctx, md)

		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}
