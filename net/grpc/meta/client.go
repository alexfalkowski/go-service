package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor returns a gRPC unary client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func UnaryClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ua, id := clientMetadata(ctx, userAgent, generator)

		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
		)
		ctx = NewOutgoingContext(ctx, md)
		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func StreamClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ua, id := clientMetadata(ctx, userAgent, generator)

		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
		)
		ctx = NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}

func clientMetadata(ctx context.Context, userAgent env.UserAgent, generator id.Generator) (Map, meta.Value, meta.Value) {
	md, ok := FromOutgoingContext(ctx)
	ua := clientUserAgent(ctx, md, userAgent)
	id := clientRequestID(ctx, generator, md)
	if !ok {
		return clientOutgoingHeaders(ua.Value(), id.Value()), ua, id
	}

	setClientOutgoingHeaders(md, ua.Value(), id.Value())

	return md, ua, id
}

func clientOutgoingHeaders(userAgent, requestID string) Map {
	// Clip caps each metadata value at one element so later appends allocate
	// instead of overwriting the neighboring value in this backing array.
	values := [...]string{userAgent, requestID}
	return Map{
		"user-agent": slices.Clip(values[0:1]),
		"request-id": slices.Clip(values[1:2]),
	}
}

func setClientOutgoingHeaders(md Map, userAgent, requestID string) {
	// Clip caps each metadata value at one element so later appends allocate
	// instead of overwriting the neighboring value in this backing array.
	values := [...]string{userAgent, requestID}
	md["user-agent"] = slices.Clip(values[0:1])
	md["request-id"] = slices.Clip(values[1:2])
}

func clientUserAgent(ctx context.Context, md metadata.MD, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}
	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	return meta.String(userAgent.String())
}

func clientRequestID(ctx context.Context, generator id.Generator, md metadata.MD) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}
	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.String(generator.Generate())
}
