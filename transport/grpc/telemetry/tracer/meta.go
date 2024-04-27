package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/grpc/meta"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/metadata"
)

// Carrier for tracer.
type Carrier struct {
	Metadata metadata.MD
}

// Get returns the value associated with the passed key.
func (s *Carrier) Get(key string) string {
	values := s.Metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set stores the key-value pair.
func (s *Carrier) Set(key string, value string) {
	s.Metadata.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (s *Carrier) Keys() []string {
	out := make([]string, 0, len(s.Metadata))
	for key := range s.Metadata {
		out = append(out, key)
	}

	return out
}

func inject(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	prop := otel.GetTextMapPropagator()

	prop.Inject(ctx, &Carrier{Metadata: md})

	return metadata.NewOutgoingContext(ctx, md)
}

func extract(ctx context.Context) context.Context {
	md := meta.ExtractIncoming(ctx)
	prop := otel.GetTextMapPropagator()

	return prop.Extract(ctx, &Carrier{Metadata: md})
}
