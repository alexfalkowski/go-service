package test

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

// WithTest stores sample info.
func WithTest(ctx context.Context, value meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, "test", value)
}

// Test retrieves sample info.
func Test(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, "test")
}
