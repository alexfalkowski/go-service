package meta

import (
	"context"
	"sync"
)

const (
	// RequestID for meta.
	RequestID = "app.request_id"
)

type contextKey string

var (
	meta            = contextKey("meta")
	mux  sync.Mutex = sync.Mutex{}
)

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
	mux.Lock()
	defer mux.Unlock()

	return attributes(ctx)[key]
}

// Attributes of meta.
func Attributes(ctx context.Context) map[string]string {
	mux.Lock()
	defer mux.Unlock()

	return attributes(ctx)
}

func attributes(ctx context.Context) map[string]string {
	m := ctx.Value(meta)
	if m == nil {
		return make(map[string]string)
	}

	return m.(map[string]string)
}
