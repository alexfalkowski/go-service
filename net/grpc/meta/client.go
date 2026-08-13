package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-service/v2/strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewClient constructs a metadata interceptor provider for a gRPC client.
func NewClient(userAgent env.UserAgent, generator id.Generator) *Client {
	return &Client{userAgent: userAgent, generator: generator}
}

// Client provides gRPC client interceptors that inject request metadata.
type Client struct {
	generator id.Generator
	userAgent env.UserAgent
}

// UnaryInterceptor returns a gRPC unary client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values plus transport metadata back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func (c *Client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ua, requestID := clientMetadata(ctx, c.userAgent, c.generator)

		ctx = meta.WithAttributes(ctx,
			meta.WithTransport(meta.Ignored("grpc")),
			meta.WithUserAgent(ua),
			meta.WithRequestID(requestID),
			meta.WithServiceMethod(meta.Ignored(fullMethod)),
		)
		ctx = NewOutgoingContext(ctx, md)
		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamInterceptor returns a gRPC stream client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values plus transport metadata back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func (c *Client) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ua, requestID := clientMetadata(ctx, c.userAgent, c.generator)

		ctx = meta.WithAttributes(ctx,
			meta.WithTransport(meta.Ignored("grpc")),
			meta.WithUserAgent(ua),
			meta.WithRequestID(requestID),
			meta.WithServiceMethod(meta.Ignored(fullMethod)),
		)
		ctx = NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}

func clientMetadata(ctx context.Context, userAgent env.UserAgent, generator id.Generator) (Map, meta.Value, meta.Value) {
	md, ok := FromOutgoingContext(ctx)
	ua := clientUserAgent(ctx, md, userAgent)
	requestID := clientRequestID(ctx, generator, md)
	wireUA := wireUserAgent(ua.Value(), userAgent.String())
	wireID := wireRequestID(requestID.Value(), generator)
	if !ok {
		return clientOutgoingHeaders(wireUA, wireID), ua, requestID
	}

	setClientOutgoingHeaders(md, wireUA, wireID)

	return md, ua, requestID
}

// wireUserAgent returns userAgent if it satisfies gRPC's printable-ASCII metadata value contract, otherwise the
// configured fallback, so a caller-supplied value outside that contract cannot fail the whole outgoing call. The
// context attribute keeps the original resolved value; only the value written to outgoing metadata is
// constrained.
func wireUserAgent(userAgent, fallback string) string {
	if header.ValidMetadataValue(userAgent) {
		return userAgent
	}

	return fallback
}

// wireRequestID returns requestID if it satisfies gRPC's printable-ASCII metadata value contract, otherwise a
// freshly generated id, so a caller-supplied value outside that contract cannot fail the whole outgoing call.
// The context attribute keeps the original resolved value; only the value written to outgoing metadata is
// constrained.
func wireRequestID(requestID string, generator id.Generator) string {
	if header.ValidMetadataValue(requestID) {
		return requestID
	}

	return generator.Generate()
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
	if ua := md.Get("user-agent"); len(ua) > 0 && !strings.IsEmpty(ua[0]) {
		return meta.String(ua[0])
	}

	return meta.String(userAgent.String())
}

func clientRequestID(ctx context.Context, generator id.Generator, md metadata.MD) meta.Value {
	if requestID := meta.RequestID(ctx); !requestID.IsEmpty() {
		return requestID
	}
	if requestID := md.Get("request-id"); len(requestID) > 0 && !strings.IsEmpty(requestID[0]) {
		return meta.String(requestID[0])
	}

	return meta.String(generator.Generate())
}
