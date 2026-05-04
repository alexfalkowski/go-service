package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/strings"
)

// NewRoundTripper constructs client-side metadata middleware for HTTP requests.
//
// The returned RoundTripper injects standard request headers and synchronizes them back into the request
// context so downstream transport wrappers (for example logging/tracing) can read consistent values.
func NewRoundTripper(userAgent env.UserAgent, generator id.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, userAgent: userAgent, generator: generator}
}

// RoundTripper wraps an underlying `http.RoundTripper` and injects request metadata.
//
// This RoundTripper is intended to be applied as an outer wrapper so that other RoundTripper middleware
// (retry/breaker/logger/token, etc.) observes the finalized headers and context values.
type RoundTripper struct {
	http.RoundTripper
	generator id.Generator
	userAgent env.UserAgent
}

// RoundTrip injects request metadata into the outbound request.
//
// It sets the "User-Agent" and "Request-Id" headers, preferring values already present in the context or
// request headers, and stores the chosen values back into the request context.
//
// Precedence rules:
//   - If the context already contains a value (meta.UserAgent/meta.RequestID), that value is used.
//   - Else, if the request header already contains a value, that value is used.
//   - Else, a default is used (userAgent parameter or a generated request id).
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := clientUserAgent(ctx, req, r.userAgent)
	requestID := clientRequestID(ctx, r.generator, req)

	clientSetRequestHeaders(req.Header, userAgent.Value(), requestID.Value())
	ctx = meta.WithAttributes(ctx,
		meta.WithUserAgent(userAgent),
		meta.WithRequestID(requestID),
	)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func clientSetRequestHeaders(header http.Header, userAgent, requestID string) {
	// One backing array avoids allocating a separate one-element slice for each header.
	values := [...]string{userAgent, requestID}
	header["User-Agent"] = values[0:1]
	header["Request-Id"] = values[1:2]
}

func clientUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); !strings.IsEmpty(ua) {
		return meta.String(ua)
	}

	return meta.String(userAgent.String())
}

func clientRequestID(ctx context.Context, generator id.Generator, req *http.Request) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}

	if id := req.Header.Get("Request-Id"); !strings.IsEmpty(id) {
		return meta.String(id)
	}

	return meta.String(generator.Generate())
}
