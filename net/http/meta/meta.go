package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

const (
	requestKey  = context.Key("request")
	responseKey = context.Key("response")
	encoderKey  = context.Key("encoder")
)

// NoPrefix is just an alias for meta.NoPrefix.
const NoPrefix = meta.NoPrefix

// Map is just an alias for meta.Map.
type Map = meta.Map

// CamelStrings is an alias for meta.CamelStrings.
func CamelStrings(ctx context.Context, prefix string) Map {
	return meta.CamelStrings(ctx, prefix)
}

// Error is an alias for meta.Error.
func Error(err error) meta.Value {
	return meta.Error(err)
}

// WithAttribute is an alias for meta.WithAttribute.
func WithAttribute(ctx context.Context, key string, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, key, value)
}

// WithRequest stores req in ctx and returns the derived context.
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, req)
}

// Request returns the stored *http.Request from ctx.
//
// Request expects WithRequest to have been called. It will panic if no request is stored in ctx.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse stores res in ctx and returns the derived context.
func WithResponse(ctx context.Context, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, res)
}

// Response returns the stored http.ResponseWriter from ctx.
//
// Response expects WithResponse to have been called. It will panic if no response is stored in ctx.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder stores enc in ctx and returns the derived context.
func WithEncoder(ctx context.Context, enc encoding.Encoder) context.Context {
	return context.WithValue(ctx, encoderKey, enc)
}

// Encoder returns the stored encoding.Encoder from ctx.
//
// Encoder expects WithEncoder to have been called. It will panic if no encoder is stored in ctx.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}
