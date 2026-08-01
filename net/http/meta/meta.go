package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

const requestResponseKey = context.Key("request-response")

const limiterKey = context.Key("limiter")

// NoPrefix is an alias for [meta.NoPrefix].
const NoPrefix = meta.NoPrefix

// Map is an alias for [meta.Map].
type Map = meta.Map

// Pair is an alias for [meta.Pair].
type Pair = meta.Pair

// CamelStrings exports all stored meta attributes as a string map with lowerCamelCased keys.
//
// The prefix parameter is prepended to each exported key (if non-empty). Attributes whose rendered value is
// empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return meta.CamelStrings(ctx, prefix)
}

// Error converts err to a [meta.Value] using err.Error().
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

type requestResponse struct {
	request  *http.Request
	response http.ResponseWriter
}

// Request returns the stored *[http.Request] from ctx.
//
// It returns nil when WithRequestResponse has not been called.
func Request(ctx context.Context) *http.Request {
	requestResponse, _ := ctx.Value(requestResponseKey).(requestResponse)

	return requestResponse.request
}

// Response returns the stored [http.ResponseWriter] from ctx.
//
// It returns nil when WithRequestResponse has not been called.
func Response(ctx context.Context) http.ResponseWriter {
	requestResponse, _ := ctx.Value(requestResponseKey).(requestResponse)

	return requestResponse.response
}

// WithRequestResponse stores the HTTP request and response writer in ctx and returns the derived context.
func WithRequestResponse(ctx context.Context, req *http.Request, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, requestResponseKey, requestResponse{request: req, response: res})
}

// RateLimiter is the minimal per-message load-control interface a streaming handler charges against.
//
// It is deliberately narrower than [github.com/alexfalkowski/go-service/v2/transport/limiter.Limiter]'s
// concrete type and defined locally: net/... packages must not import transport/... (see AGENTS.md).
// *[github.com/alexfalkowski/go-service/v2/transport/limiter.Limiter] already satisfies this interface
// structurally via its Take method, so [transport/http/limiter.Handler.ServeHTTP] can pass it to
// WithLimiter without an adapter.
type RateLimiter interface {
	// Take attempts to take a token for the key derived from ctx. It reports whether the request is
	// allowed and any error from the underlying store; the header value returned alongside the same
	// decision by the concrete limiter is intentionally not part of this interface, since RateLimit
	// headers cannot be re-sent once a streaming response is committed.
	Take(ctx context.Context) (bool, string, error)
}

// WithLimiter stores limiter in ctx and returns the derived context.
//
// This is a separate context key from WithRequestResponse because a limiter is route-specific and usually
// absent. limiter is typically the same limiter that already charged one token for the request at stream open
// (see [transport/http/limiter.Handler.ServeHTTP]); storing it here lets a streaming handler, constructed
// deeper in the call chain, charge additional per-message tokens against the same limiter. limiter may be
// nil, meaning no limiter is configured for the request.
func WithLimiter(ctx context.Context, limiter RateLimiter) context.Context {
	return context.WithValue(ctx, limiterKey, limiter)
}

// Limiter returns the stored [RateLimiter] from ctx, or nil if WithLimiter was never called or was called
// with a nil limiter.
//
// Like Request and Response, Limiter returns nil on a missing value. A limiter is optional per-route
// configuration, so nil ("no limiter configured") is the common case rather than a caller error.
func Limiter(ctx context.Context) RateLimiter {
	limiter, _ := ctx.Value(limiterKey).(RateLimiter)
	return limiter
}
