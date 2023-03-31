package otel

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func extract(ctx context.Context, h message.Headers) context.Context {
	prop := otel.GetTextMapPropagator()

	return prop.Extract(ctx, headerCarrier(h))
}

func inject(ctx context.Context, h message.Headers) {
	prop := otel.GetTextMapPropagator()

	prop.Inject(ctx, headerCarrier(h))
}

type headerCarrier message.Headers

var _ propagation.TextMapCarrier = &headerCarrier{}

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return hc[key]
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	hc[key] = value
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}

	return keys
}
