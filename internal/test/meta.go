package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
)

// WithTest creates a test metadata pair for meta.WithAttributes.
//
// The pair stores value under the `test` key used by shared test helpers.
func WithTest(value meta.Value) meta.Pair {
	return meta.NewPair("test", value)
}

// Test retrieves the metadata attribute stored under the `test` key.
//
// If no value is present, this returns the zero-value meta.Value.
func Test(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, "test")
}
