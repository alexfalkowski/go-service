package meta

import (
	"context"
)

type contextKey string

var meta = contextKey("meta")

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key, value string) context.Context {
	attr := attributes(ctx)
	attr[key] = value

	return context.WithValue(ctx, meta, attr)
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) string {
	return attributes(ctx)[key]
}

// Attributes of meta.
func Attributes(ctx context.Context) map[string]string {
	return attributes(ctx)
}

func attributes(ctx context.Context) map[string]string {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]string)
	}

	return m.(map[string]string)
}
