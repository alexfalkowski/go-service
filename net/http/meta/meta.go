package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

var (
	// CamelStrings is an alias for meta.CamelStrings.
	CamelStrings = meta.CamelStrings

	// Error is an alias for meta.Error.
	Error = meta.Error

	// WithAttribute is an alias for meta.WithAttribute.
	WithAttribute = meta.WithAttribute
)

type contextKey string

const (
	requestKey  = contextKey("request")
	responseKey = contextKey("response")
	encoderKey  = contextKey("encoder")
)

// WithRequest for http.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

// Request for http.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse for http.
func WithResponse(ctx context.Context, r http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, r)
}

// Response for http.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder for http.
func WithEncoder(ctx context.Context, e encoding.Encoder) context.Context {
	return context.WithValue(ctx, encoderKey, e)
}

// Encoder for rpc.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}
