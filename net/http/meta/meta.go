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
