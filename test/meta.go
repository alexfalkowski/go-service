package test

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/meta"
)

// WithTest stores sample info.
func WithTest(ctx context.Context, value fmt.Stringer) context.Context {
	return meta.WithAttribute(ctx, "test", value)
}

// Test retrieves sample info.
func Test(ctx context.Context) fmt.Stringer {
	return meta.Attribute(ctx, "test")
}
