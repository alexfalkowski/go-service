package tracer

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for tracer.
func UnaryServerInterceptor(trace *tracer.Tracer) grpc.UnaryServerInterceptor {
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

		ctx, span := trace.StartServer(ctx, operationName(info.FullMethod), attrs...)
		defer span.End()

		resp, err := handler(ctx, req)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return resp, err
	}
}

// StreamServerInterceptor for tracer.
func StreamServerInterceptor(trace *tracer.Tracer) grpc.StreamServerInterceptor {
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

		ctx, span := trace.StartServer(ctx, operationName(info.FullMethod), attrs...)
		defer span.End()

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
func UnaryClientInterceptor(trace *tracer.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := trace.StartClient(ctx, operationName(conn.Target()+fullMethod), attrs...)
		defer span.End()

		ctx = inject(ctx)

		err := invoker(ctx, fullMethod, req, resp, conn, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return err
	}
}

// StreamClientInterceptor for tracer.
func StreamClientInterceptor(trace *tracer.Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		method := path.Base(fullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}

		ctx, span := trace.StartClient(ctx, operationName(conn.Target()+fullMethod), attrs...)
		defer span.End()

		ctx = inject(ctx)

		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(status.Code(err))))

		return stream, err
	}
}

func operationName(name string) string {
	return tracer.OperationName("grpc", "get "+name)
}
