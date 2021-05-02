package grpc

import (
	"context"
	"strings"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/google/uuid"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type metadataTextMap metadata.MD

// Set is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) Set(key, val string) {
	key = strings.ToLower(key)

	m[key] = append(m[key], val)
}

// ForeachKey is a opentracing.TextMapReader interface that extracts values.
func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if err := callback(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func metaUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := extractIncoming(ctx)

		var requestID string

		if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
			requestID = mdRequestID[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithAttribute(ctx, meta.RequestID, requestID)
		header := metadata.Pairs("request-id", requestID)

		if err := grpc.SendHeader(ctx, header); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func metaStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := extractIncoming(ctx)

		var requestID string

		if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
			requestID = mdRequestID[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithAttribute(ctx, meta.RequestID, requestID)
		header := metadata.Pairs("request-id", requestID)

		if err := grpc.SendHeader(ctx, header); err != nil {
			return err
		}

		wrappedStream := grpcMiddleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, stream)
	}
}

func metaUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := meta.Attribute(ctx, meta.RequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithAttribute(ctx, meta.RequestID, requestID)
		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

func metaStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requestID := meta.Attribute(ctx, meta.RequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithAttribute(ctx, meta.RequestID, requestID)
		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}
