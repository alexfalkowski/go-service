package tracer

import (
	"path"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
)

// Tracer is an alias for tracer.Tracer.
type Tracer = tracer.Tracer

// UnaryServerInterceptor for tracer.
func UnaryServerInterceptor(trace *Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		ctx = extract(ctx)
		method := path.Base(info.FullMethod)

		ctx, span := trace.StartServer(ctx, operationName(info.FullMethod),
			attributes.RPCSystemGRPC,
			attributes.RPCService(p),
			attributes.RPCMethod(method))
		defer span.End()

		resp, err := handler(ctx, req)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)

		span.SetAttributes(attributes.GRPCStatusCode(int64(status.Code(err))))

		return resp, err
	}
}

// StreamServerInterceptor for tracer.
func StreamServerInterceptor(trace *Tracer) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		ctx := extract(stream.Context())
		method := path.Base(info.FullMethod)

		ctx, span := trace.StartServer(ctx, operationName(info.FullMethod),
			attributes.RPCSystemGRPC,
			attributes.RPCService(p),
			attributes.RPCMethod(method))
		defer span.End()

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(attributes.GRPCStatusCode(int64(status.Code(err))))

		return err
	}
}

// UnaryClientInterceptor for tracer.
func UnaryClientInterceptor(trace *Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		p := fullMethod[1:]
		if strings.IsObservable(p) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		method := path.Base(fullMethod)

		ctx, span := trace.StartClient(ctx, operationName(conn.Target()+fullMethod),
			attributes.RPCSystemGRPC,
			attributes.RPCService(p),
			attributes.RPCMethod(method))
		defer span.End()

		ctx = inject(ctx)

		err := invoker(ctx, fullMethod, req, resp, conn, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(attributes.GRPCStatusCode(int64(status.Code(err))))

		return err
	}
}

// StreamClientInterceptor for tracer.
func StreamClientInterceptor(trace *Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		p := fullMethod[1:]
		if strings.IsObservable(p) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		method := path.Base(fullMethod)

		ctx, span := trace.StartClient(ctx, operationName(conn.Target()+fullMethod),
			attributes.RPCSystemGRPC,
			attributes.RPCService(p),
			attributes.RPCMethod(method))
		defer span.End()

		ctx = inject(ctx)

		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)

		tracer.Error(err, span)
		tracer.Meta(ctx, span)
		span.SetAttributes(attributes.GRPCStatusCode(int64(status.Code(err))))

		return stream, err
	}
}

func operationName(name string) string {
	return tracer.OperationName("grpc", "get "+name)
}
