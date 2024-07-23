package context

import (
	"context"
	"net/http"
)

type contextKey string

var (
	requestKey  = contextKey("request")
	responseKey = contextKey("response")
)

// WithRequest for rpc.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

// Request for rpc.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse for rpc.
func WithResponse(ctx context.Context, r http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, r)
}

// Response for rpc.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}
