package context

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
)

type (
	// Context is an alias of Context.
	Context = context.Context

	contextKey string
)

const (
	requestKey  = contextKey("request")
	responseKey = contextKey("response")
	encoderKey  = contextKey("encoder")
)

// WithRequest for http.
func WithRequest(ctx Context, r *http.Request) Context {
	return context.WithValue(ctx, requestKey, r)
}

// Request for http.
func Request(ctx Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse for http.
func WithResponse(ctx Context, r http.ResponseWriter) Context {
	return context.WithValue(ctx, responseKey, r)
}

// Response for http.
func Response(ctx Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder for http.
func WithEncoder(ctx Context, e encoding.Encoder) Context {
	return context.WithValue(ctx, encoderKey, e)
}

// Encoder for rpc.
func Encoder(ctx Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}
