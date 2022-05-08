package opentracing

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	sstrings "github.com/alexfalkowski/go-service/strings"
	stime "github.com/alexfalkowski/go-service/time"
	sopentracing "github.com/alexfalkowski/go-service/trace/opentracing"
	sgmeta "github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/version"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	grpcService         = "grpc.service"
	grpcMethod          = "grpc.method"
	grpcCode            = "grpc.code"
	grpcDuration        = "grpc.duration"
	grpcStartTime       = "grpc.start_time"
	grpcRequestDeadline = "grpc.request.deadline"
	component           = "component"
	grpcComponent       = "grpc"
)

// TracerParams for opentracing.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *sopentracing.Config
	Version   version.Version
}

// NewTracer for opentracing.
func NewTracer(params TracerParams) (Tracer, error) {
	return sopentracing.NewTracer(sopentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "grpc", Config: params.Config, Version: params.Version})
}

// Tracer for opentracing.
type Tracer opentracing.Tracer

// UnaryServerInterceptor for opentracing.
func UnaryServerInterceptor(tracer Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if sstrings.IsHealth(service) {
			return handler(ctx, req)
		}

		start := time.Now().UTC()
		md := sgmeta.ExtractIncoming(ctx)
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

		span := tracer.StartSpan(operationName(info.FullMethod), opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = opentracing.ContextWithSpan(ctx, span)
		resp, err := handler(ctx, req)

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcDuration, stime.ToMilliseconds(time.Since(start)))
		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		return resp, err
	}
}

// StreamServerInterceptor for opentracing.
func StreamServerInterceptor(tracer Tracer) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if sstrings.IsHealth(service) {
			return handler(srv, stream)
		}

		start := time.Now().UTC()
		ctx := stream.Context()
		md := sgmeta.ExtractIncoming(ctx)
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

		span := tracer.StartSpan(operationName(info.FullMethod), opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = opentracing.ContextWithSpan(ctx, span)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcDuration, stime.ToMilliseconds(time.Since(start)))
		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		return err
	}
}

// UnaryClientInterceptor for opentracing.
func UnaryClientInterceptor(tracer Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if sstrings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		spanOpts := []opentracing.StartSpanOption{
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName(fullMethod), spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := sgmeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			return err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, fullMethod, req, resp, cc, opts...)

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcDuration, stime.ToMilliseconds(time.Since(start)))
		span.SetTag(grpcCode, status.Code(err))

		if err != nil {
			setError(span, err)
		}

		return err
	}
}

// StreamClientInterceptor for opentracing.
func StreamClientInterceptor(tracer Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if sstrings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now().UTC()
		method := path.Base(fullMethod)
		spanOpts := []opentracing.StartSpanOption{
			opentracing.Tag{Key: grpcStartTime, Value: start.Format(time.RFC3339)},
			opentracing.Tag{Key: grpcService, Value: service},
			opentracing.Tag{Key: grpcMethod, Value: method},
			opentracing.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName(fullMethod), spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := sgmeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			return nil, err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcDuration, stime.ToMilliseconds(time.Since(start)))
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
	tags := tags.Extract(ctx)
	for k, v := range tags.Values() {
		if err, ok := v.(error); ok {
			span.LogKV(k, err.Error())
		} else {
			span.SetTag(k, v)
		}
	}
}

func operationName(fullMethod string) string {
	return fmt.Sprintf("get %s", fullMethod)
}
