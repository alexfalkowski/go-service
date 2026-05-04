package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

const contentKey = context.Key("content")

// NoPrefix is an alias for meta.NoPrefix.
const NoPrefix = meta.NoPrefix

// Map is an alias for meta.Map.
type Map = meta.Map

// Pair is an alias for meta.Pair.
type Pair = meta.Pair

type content struct {
	request  *http.Request
	response http.ResponseWriter
	encoder  encoding.Encoder
}

// CamelStrings exports all stored meta attributes as a string map with lowerCamelCased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty). Attributes whose rendered value is
// empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return meta.CamelStrings(ctx, prefix)
}

// Error converts err to a meta.Value using err.Error().
func Error(err error) meta.Value {
	return meta.Error(err)
}

// NewPair creates one metadata key/value pair for batched storage updates.
func NewPair(key string, value meta.Value) Pair {
	return meta.NewPair(key, value)
}

// WithAttributes stores all provided metadata pairs on ctx.
func WithAttributes(ctx context.Context, pairs ...Pair) context.Context {
	return meta.WithAttributes(ctx, pairs...)
}

// Request returns the stored *http.Request from ctx.
//
// Panics: Request expects WithContent to have been called. It will panic if no content metadata is stored
// in ctx or if the stored value is invalid.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(contentKey).(content).request
}

// Response returns the stored http.ResponseWriter from ctx.
//
// Panics: Response expects WithContent to have been called. It will panic if no content metadata is stored
// in ctx or if the stored value is invalid.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(contentKey).(content).response
}

// Encoder returns the stored encoding.Encoder from ctx.
//
// Panics: Encoder expects WithContent to have been called. It will panic if no content metadata is stored
// in ctx or if the stored value is invalid.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(contentKey).(content).encoder
}

// WithContent stores HTTP content metadata in ctx and returns the derived context.
//
// The encoder may be nil when a handler only needs request and response access.
func WithContent(ctx context.Context, req *http.Request, res http.ResponseWriter, enc encoding.Encoder) context.Context {
	return context.WithValue(ctx, contentKey, content{request: req, response: res, encoder: enc})
}
