package tracer

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for tracer.
func UnaryServerInterceptor(t trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(ctx, req)
		}

		ctx = extract(ctx)

		method := path.Base(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := t.Start(trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)), operationName(info.FullMethod),
			trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...))
		defer span.End()

		ctx = tracer.WithTraceID(ctx, span)
		resp, err := handler(ctx, req)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return resp, err
	}
}

// StreamServerInterceptor for tracer.
func StreamServerInterceptor(t trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(srv, stream)
		}

		ctx := extract(stream.Context())

		method := path.Base(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := t.Start(trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)), operationName(info.FullMethod),
			trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...))
		defer span.End()

		ctx = tracer.WithTraceID(ctx, span)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return err
	}
}

// UnaryClientInterceptor for tracer.
func UnaryClientInterceptor(t trace.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := t.Start(ctx, operationName(cc.Target()+fullMethod), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
		defer span.End()

		ctx = tracer.WithTraceID(ctx, span)
		ctx = inject(ctx)

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return err
	}
}

// StreamClientInterceptor for tracer.
func StreamClientInterceptor(t trace.Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := t.Start(ctx, operationName(cc.Target()+fullMethod), trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))
		defer span.End()

		ctx = tracer.WithTraceID(ctx, span)
		ctx = inject(ctx)

		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return stream, err
	}
}

func operationName(name string) string {
	return tracer.OperationName("grpc", "get "+name)
}
