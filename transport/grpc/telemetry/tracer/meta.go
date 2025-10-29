package tracer

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
)

// NewCarrier for tracer.
func NewCarrier(meta meta.Map) *Carrier {
	return &Carrier{meta: meta}
}

// Carrier for tracer.
type Carrier struct {
	meta meta.Map
}

// Get returns the value associated with the passed key.
func (s *Carrier) Get(key string) string {
	values := s.meta.Get(key)
	if len(values) == 0 {
		return strings.Empty
	}

	return values[0]
}

// Set stores the key-value pair.
func (s *Carrier) Set(key string, value string) {
	s.meta.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (s *Carrier) Keys() []string {
	out := make([]string, len(s.meta))
	cnt := 0

	for key := range s.meta {
		out[cnt] = key
		cnt++
	}

	return out
}

func inject(ctx context.Context) context.Context {
	md := meta.ExtractOutgoing(ctx)
	telemetry.Inject(ctx, NewCarrier(md))

	return meta.NewOutgoingContext(ctx, md)
}

func extract(ctx context.Context) context.Context {
	md := meta.ExtractIncoming(ctx)

	return telemetry.Extract(ctx, NewCarrier(md))
}
