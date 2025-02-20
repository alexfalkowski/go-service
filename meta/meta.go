package meta

import (
	"context"

	"github.com/samber/lo"
)

type contextKey string

const meta = contextKey("meta")

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key string, value Value) context.Context {
	return context.WithValue(ctx, meta, attributes(ctx).Add(key, value))
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) Value {
	return attributes(ctx).Get(key)
}

// Map for meta.
type Map map[string]string

// SnakeStrings for meta.
func SnakeStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, lo.SnakeCase)
}

// CamelStrings for meta.
func CamelStrings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, lo.CamelCase)
}

// Strings for meta.
func Strings(ctx context.Context, prefix string) Map {
	return attributes(ctx).Strings(prefix, func(s string) string { return s })
}
