package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// WithTest stores sample info.
func WithTest(ctx context.Context, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, "test", value)
}

// Test retrieves sample info.
func Test(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, "test")
}
