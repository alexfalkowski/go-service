package opentracing

import (
	"context"
	"path"

	grpcMeta "github.com/alexfalkowski/go-service/pkg/grpc/meta"
	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	grpcRequest         = "grpc.request"
	grpcResponse        = "grpc.response"
	grpcService         = "grpc.service"
	grpcMethod          = "grpc.method"
	grpcCode            = "grpc.code"
	grpcDuration        = "grpc.duration"
	grpcStartTime       = "grpc.start_time"
	grpcRequestDeadline = "grpc.request.deadline"
	component           = "component"
	grpcComponent       = "grpc"
	healthService       = "grpc.health.v1.Health"
)

// UnaryServerInterceptor for opentracing.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(ctx, req)
		}

		start := time.Now().UTC()
		tracer := opentracing.GlobalTracer()
		md := grpcMeta.ExtractIncoming(ctx)
		method := path.Base(info.FullMethod)
		traceCtx, _ := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		opts := []opentracing.StartSpanOption{
			ext.RPCServerOption(traceCtx),
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: grpcRequest, Value: encode(req)},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCServer,
		}

		for k, v := range meta.Attributes(ctx) {
			opts = append(opts, opentracing.Tag{Key: k, Value: v})
		}

		span := tracer.StartSpan(info.FullMethod, opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = opentracing.ContextWithSpan(ctx, span)
		resp, err := handler(ctx, req)

		span.SetTag(grpcDuration, time.ToMilliseconds(time.Since(start)))
		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		if resp != nil {
			span.SetTag(grpcResponse, encode(resp))
		}

		return resp, err
	}
}

// StreamServerInterceptor for opentracing.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if service == healthService {
			return handler(srv, stream)
		}

		start := time.Now().UTC()
		ctx := stream.Context()
		tracer := opentracing.GlobalTracer()
		md := grpcMeta.ExtractIncoming(ctx)
		method := path.Base(info.FullMethod)
		traceCtx, _ := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		opts := []opentracing.StartSpanOption{
			ext.RPCServerOption(traceCtx),
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCServer,
		}

		for k, v := range meta.Attributes(ctx) {
			opts = append(opts, opentracing.Tag{Key: k, Value: v})
		}

		span := tracer.StartSpan(info.FullMethod, opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = opentracing.ContextWithSpan(ctx, span)

		wrappedStream := grpcMiddleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)

		span.SetTag(grpcDuration, time.ToMilliseconds(time.Since(start)))
		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		return err
	}
}

// UnaryClientInterceptor for opentracing.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if service == healthService {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		tracer := opentracing.GlobalTracer()
		spanOpts := []opentracing.StartSpanOption{
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: grpcRequest, Value: encode(req)},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		for k, v := range meta.Attributes(ctx) {
			spanOpts = append(spanOpts, opentracing.Tag{Key: k, Value: v})
		}

		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, fullMethod, spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := grpcMeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			return err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, fullMethod, req, resp, cc, opts...)

		span.SetTag(grpcDuration, time.ToMilliseconds(time.Since(start)))
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		if resp != nil {
			span.SetTag(grpcResponse, encode(resp))
		}

		return err
	}
}

// StreamClientInterceptor for opentracing.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if service == healthService {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		tracer := opentracing.GlobalTracer()
		spanOpts := []opentracing.StartSpanOption{
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		for k, v := range meta.Attributes(ctx) {
			spanOpts = append(spanOpts, opentracing.Tag{Key: k, Value: v})
		}

		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, fullMethod, spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := grpcMeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			return nil, err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)

		span.SetTag(grpcDuration, time.ToMilliseconds(time.Since(start)))
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		return stream, err
	}
}

func setError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
}

func addTags(ctx context.Context, span opentracing.Span) {
	tags := grpcTags.Extract(ctx)
	for k, v := range tags.Values() {
		if err, ok := v.(error); ok {
			span.LogKV(k, err.Error())
		} else {
			span.SetTag(k, v)
		}
	}
}
