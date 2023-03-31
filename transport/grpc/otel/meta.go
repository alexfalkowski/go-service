package otel

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func inject(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	prop := otel.GetTextMapPropagator()

	prop.Inject(ctx, &metadataCarrier{metadata: &md})

	return metadata.NewOutgoingContext(ctx, md)
}

func extract(ctx context.Context) context.Context {
	md := meta.ExtractIncoming(ctx)
	prop := otel.GetTextMapPropagator()

	return prop.Extract(ctx, &metadataCarrier{metadata: &md})
}

type metadataCarrier struct {
	metadata *metadata.MD
}

var _ propagation.TextMapCarrier = &metadataCarrier{}

// Get returns the value associated with the passed key.
func (s *metadataCarrier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set stores the key-value pair.
func (s *metadataCarrier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (s *metadataCarrier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}

	return out
}
