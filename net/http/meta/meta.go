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

// NoPrefix is an alias for meta.NoPrefix.
const NoPrefix = meta.NoPrefix

// Map is an alias for meta.Map.
type Map = meta.Map

// CamelStrings exports all stored meta attributes as a string map with lowerCamelCased keys.
//
// This is a thin wrapper around meta.CamelStrings. The prefix parameter is prepended to each exported key
// (if non-empty). Attributes whose rendered value is empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return meta.CamelStrings(ctx, prefix)
}

// Error converts err to a meta.Value using err.Error().
//
// This is a thin wrapper around meta.Error.
func Error(err error) meta.Value {
	return meta.Error(err)
}

// WithAttribute stores an arbitrary meta attribute on ctx.
//
// This is a thin wrapper around meta.WithAttribute.
func WithAttribute(ctx context.Context, key string, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, key, value)
}

// WithRequest stores req in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the request available to
// downstream handlers via Request(ctx).
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, req)
}

// Request returns the stored *http.Request from ctx.
//
// Panics: Request expects WithRequest to have been called. It will panic if no request is stored in ctx
// or if the stored value is not a *http.Request.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse stores res in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the response writer available
// to downstream handlers via Response(ctx).
func WithResponse(ctx context.Context, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, res)
}

// Response returns the stored http.ResponseWriter from ctx.
//
// Panics: Response expects WithResponse to have been called. It will panic if no response writer is stored
// in ctx or if the stored value is not an http.ResponseWriter.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder stores enc in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the negotiated encoder
// (selected from Content-Type) available to downstream handlers via Encoder(ctx).
func WithEncoder(ctx context.Context, enc encoding.Encoder) context.Context {
	return context.WithValue(ctx, encoderKey, enc)
}

// Encoder returns the stored encoding.Encoder from ctx.
//
// Panics: Encoder expects WithEncoder to have been called. It will panic if no encoder is stored in ctx
// or if the stored value is not an encoding.Encoder.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}
