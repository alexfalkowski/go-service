package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// WithTest stores a metadata attribute named `test` on the context.
func WithTest(ctx context.Context, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, "test", value)
}

// Test retrieves the metadata attribute stored under the `test` key.
func Test(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, "test")
}
