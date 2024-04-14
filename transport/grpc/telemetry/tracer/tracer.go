package tracer

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/meta"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for tracer.
func UnaryServerInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		ctx = extract(ctx)

		method := path.Base(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
			operationName(info.FullMethod),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

		resp, err := handler(ctx, req)
		if err != nil {
			span.SetStatus(codes.Error, status.Code(err).String())
			span.RecordError(err)
		}

		for k, v := range meta.Strings(ctx) {
			span.SetAttributes(attribute.Key(k).String(v))
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return resp, err
	}
}

// StreamServerInterceptor for tracer.
func StreamServerInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		ctx := extract(stream.Context())

		method := path.Base(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
			operationName(info.FullMethod),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)
		if err != nil {
			span.SetStatus(codes.Error, status.Code(err).String())
			span.RecordError(err)
		}

		for k, v := range meta.Strings(ctx) {
			span.SetAttributes(attribute.Key(k).String(v))
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return err
	}
}

// UnaryClientInterceptor for tracer.
func UnaryClientInterceptor(tracer trace.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := tracer.Start(
			ctx,
			operationName(cc.Target()+fullMethod),
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
		ctx = inject(ctx)

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		if err != nil {
			span.SetStatus(codes.Error, status.Code(err).String())
			span.RecordError(err)
		}

		for k, v := range meta.Strings(ctx) {
			span.SetAttributes(attribute.Key(k).String(v))
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return err
	}
}

// StreamClientInterceptor for tracer.
func StreamClientInterceptor(tracer trace.Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := tracer.Start(
			ctx,
			operationName(cc.Target()+fullMethod),
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		ctx = tm.WithTraceID(ctx, meta.ToValuer(span.SpanContext().TraceID()))
		ctx = inject(ctx)

		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		if err != nil {
			span.SetStatus(codes.Error, status.Code(err).String())
			span.RecordError(err)
		}

		for k, v := range meta.Strings(ctx) {
			span.SetAttributes(attribute.Key(k).String(v))
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return stream, err
	}
}

func operationName(fullMethod string) string {
	return "get " + fullMethod
}
