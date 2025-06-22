package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

type contextKey string

const (
	requestKey  = contextKey("request")
	responseKey = contextKey("response")
	encoderKey  = contextKey("encoder")
)

// Map is just an alias for meta.Map.
type Map = meta.Map

var (
	// CamelStrings is an alias for meta.CamelStrings.
	CamelStrings = meta.CamelStrings

	// Error is an alias for meta.Error.
	Error = meta.Error

	// WithAttribute is an alias for meta.WithAttribute.
	WithAttribute = meta.WithAttribute
)

// WithRequest for http.
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, req)
}

// Request for http.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse for http.
func WithResponse(ctx context.Context, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, res)
}

// Response for http.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder for http.
func WithEncoder(ctx context.Context, enc encoding.Encoder) context.Context {
	return context.WithValue(ctx, encoderKey, enc)
}

// Encoder for rpc.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}
