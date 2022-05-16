package opentracing

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	gmeta "github.com/alexfalkowski/go-service/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/version"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	grpcService         = "grpc.service"
	grpcMethod          = "grpc.method"
	grpcCode            = "grpc.code"
	grpcRequestDeadline = "grpc.request.deadline"
	component           = "component"
	grpcComponent       = "grpc"
)

// TracerParams for otr.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for otr.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "grpc", Config: params.Config, Version: params.Version})
}

// Tracer for otr.
type Tracer otr.Tracer

// UnaryServerInterceptor for otr.
func UnaryServerInterceptor(tracer Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		md := gmeta.ExtractIncoming(ctx)
		method := path.Base(info.FullMethod)
		traceCtx, _ := tracer.Extract(otr.HTTPHeaders, metadataTextMap(md))
		opts := []otr.StartSpanOption{
			ext.RPCServerOption(traceCtx),
			otr.Tag{Key: grpcService, Value: service},
			otr.Tag{Key: grpcMethod, Value: method},
			otr.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCServer,
		}

		span := tracer.StartSpan(operationName(info.FullMethod), opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = otr.ContextWithSpan(ctx, span)

		resp, err := handler(ctx, req)
		if err != nil {
			opentracing.SetError(span, err)
		}

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		return resp, err
	}
}

// StreamServerInterceptor for otr.
func StreamServerInterceptor(tracer Tracer) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		md := gmeta.ExtractIncoming(ctx)
		method := path.Base(info.FullMethod)
		traceCtx, _ := tracer.Extract(otr.HTTPHeaders, metadataTextMap(md))
		opts := []otr.StartSpanOption{
			ext.RPCServerOption(traceCtx),
			otr.Tag{Key: grpcService, Value: service},
			otr.Tag{Key: grpcMethod, Value: method},
			otr.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCServer,
		}

		span := tracer.StartSpan(operationName(info.FullMethod), opts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		ctx = otr.ContextWithSpan(ctx, span)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		err := handler(srv, wrappedStream)
		if err != nil {
			opentracing.SetError(span, err)
		}

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		addTags(ctx, span)
		span.SetTag(grpcCode, status.Code(err))

		return err
	}
}

// UnaryClientInterceptor for otr.
func UnaryClientInterceptor(tracer Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		method := path.Base(fullMethod)
		spanOpts := []otr.StartSpanOption{
			otr.Tag{Key: grpcService, Value: service},
			otr.Tag{Key: grpcMethod, Value: method},
			otr.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		span, ctx := otr.StartSpanFromContextWithTracer(ctx, tracer, operationName(fullMethod), spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := gmeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), otr.HTTPHeaders, metadataTextMap(md)); err != nil {
			return err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		if err != nil {
			opentracing.SetError(span, err)
		}

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcCode, status.Code(err))

		return err
	}
}

// StreamClientInterceptor for otr.
func StreamClientInterceptor(tracer Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		method := path.Base(fullMethod)
		spanOpts := []otr.StartSpanOption{
			otr.Tag{Key: grpcService, Value: service},
			otr.Tag{Key: grpcMethod, Value: method},
			otr.Tag{Key: component, Value: grpcComponent},
			ext.SpanKindRPCClient,
		}

		span, ctx := otr.StartSpanFromContextWithTracer(ctx, tracer, operationName(fullMethod), spanOpts...)
		defer span.Finish()

		if d, ok := ctx.Deadline(); ok {
			span.SetTag(grpcRequestDeadline, d.UTC().Format(time.RFC3339))
		}

		md := gmeta.ExtractOutgoing(ctx)
		if err := tracer.Inject(span.Context(), otr.HTTPHeaders, metadataTextMap(md)); err != nil {
			return nil, err
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		if err != nil {
			opentracing.SetError(span, err)
		}

		for k, v := range meta.Attributes(ctx) {
			span.SetTag(k, v)
		}

		span.SetTag(grpcCode, status.Code(err))

		return stream, err
	}
}

func addTags(ctx context.Context, span otr.Span) {
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
