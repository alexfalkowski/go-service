package meta

import (
	"context"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	tmeta "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/version"
	"github.com/google/uuid"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(version version.Version) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := ExtractIncoming(ctx)
		ctx = meta.WithVersion(ctx, string(version))
		userAgent := extractUserAgent(ctx, md)
		ctx = tmeta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = tmeta.WithRequestID(ctx, requestID)

		remoteAddress := extractRemoteAddress(ctx, md)
		ctx = tmeta.WithRemoteAddress(ctx, remoteAddress)

		headers := metadata.Pairs("version", string(version), "request-id", requestID, "ua", userAgent)
		if err := grpc.SendHeader(ctx, headers); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(version version.Version) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := ExtractIncoming(ctx)
		ctx = meta.WithVersion(ctx, string(version))

		userAgent := extractUserAgent(ctx, md)
		ctx = tmeta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = tmeta.WithRequestID(ctx, requestID)
		ctx = tmeta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, md))

		headers := metadata.Pairs("version", string(version), "request-id", requestID, "ua", userAgent)
		if err := grpc.SendHeader(ctx, headers); err != nil {
			return err
		}

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, stream)
	}
}

// UnaryClientInterceptor for meta.
func UnaryClientInterceptor(userAgent string, version version.Version) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)
		ctx = meta.WithVersion(ctx, string(version))

		ua := extractUserAgent(ctx, md)
		if ua == "" {
			ua = userAgent
		}

		ctx = tmeta.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = tmeta.WithRequestID(ctx, requestID)
		ctx = tmeta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, md))
		ctx = metadata.AppendToOutgoingContext(ctx, "version", string(version), "request-id", requestID, "ua", ua)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent string, version version.Version) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)
		ctx = meta.WithVersion(ctx, string(version))

		ua := extractUserAgent(ctx, md)
		if ua == "" {
			ua = userAgent
		}

		ctx = tmeta.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = tmeta.WithRequestID(ctx, requestID)
		ctx = tmeta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, md))
		ctx = metadata.AppendToOutgoingContext(ctx, "version", string(version), "request-id", requestID, "ua", ua)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}

func extractUserAgent(ctx context.Context, md metadata.MD) string {
	if mdUserAgent := md.Get("ua"); len(mdUserAgent) > 0 {
		return mdUserAgent[0]
	}

	return tmeta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, md metadata.MD) string {
	if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
		return mdRequestID[0]
	}

	return tmeta.RequestID(ctx)
}

func extractRemoteAddress(ctx context.Context, md metadata.MD) string {
	if mdfForwardedFor := md.Get("forwarded-for"); len(mdfForwardedFor) > 0 {
		return strings.Split(mdfForwardedFor[0], ",")[0]
	}

	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return tmeta.RemoteAddress(ctx)
}
