package rpc

import (
	"context"
	"net/http"
)

type contextKey string

// WithRequest for rpc.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, contextKey("request"), r)
}

// Request for rpc.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(contextKey("request")).(*http.Request)
}

// WithResponse for rpc.
func WithResponse(ctx context.Context, r http.ResponseWriter) context.Context {
	return context.WithValue(ctx, contextKey("response"), r)
}

// Response for rpc.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(contextKey("response")).(http.ResponseWriter)
}
