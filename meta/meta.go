package meta

import (
	"context"
	"sync"
)

type contextKey string

// nolint:gochecknoglobals
var (
	meta = contextKey("meta")
	mux  = sync.RWMutex{}
)

const (
	versionKey = "version"
)

// WithVersion for meta.
func WithVersion(ctx context.Context, version string) context.Context {
	return WithAttribute(ctx, versionKey, version)
}

// Version for meta.
func Version(ctx context.Context) string {
	return Attribute(ctx, versionKey)
}

// WithAttribute to meta.
func WithAttribute(ctx context.Context, key, value string) context.Context {
	mux.Lock()
	defer mux.Unlock()

	attr := attributes(ctx)
	attr[key] = value

	return context.WithValue(ctx, meta, attr)
}

// Attribute of meta.
func Attribute(ctx context.Context, key string) string {
	mux.RLock()
	defer mux.RUnlock()

	return attributes(ctx)[key]
}

// Attributes of meta.
func Attributes(ctx context.Context) map[string]string {
	mux.RLock()
	defer mux.RUnlock()

	return attributes(ctx)
}

// nolint:forcetypeassert
func attributes(ctx context.Context) map[string]string {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]string)
	}

	return m.(map[string]string)
}
