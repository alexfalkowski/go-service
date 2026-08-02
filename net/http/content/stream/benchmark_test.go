package stream_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

var benchmarkMedia stream.Media

// BenchmarkNewFromMediaNDJSON tracks the exact NDJSON media-type fast path used on streaming request paths.
func BenchmarkNewFromMediaNDJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromMedia(media.NDJSON)
	}
}

// BenchmarkNewFromMediaWithParameters tracks the parser path for parameterized NDJSON media types.
func BenchmarkNewFromMediaWithParameters(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromMedia(media.NDJSON + "; profile=test")
	}
}

// BenchmarkNewFromAcceptNDJSON tracks response media negotiation for a common NDJSON stream.
func BenchmarkNewFromAcceptNDJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), http.MethodGet, "/hello", nil)
	req.Header.Set(http.AcceptKey, media.NDJSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromAccept(req)
	}
}

// BenchmarkNewFromAcceptWithRanges tracks Accept-list parsing and wildcard matching for a browser-style request.
func BenchmarkNewFromAcceptWithRanges(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), http.MethodGet, "/hello", nil)
	req.Header.Set(http.AcceptKey, "text/html,application/xhtml+xml,*/*;q=0.8")

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromAccept(req)
	}
}

// BenchmarkNewFromContentTypeNDJSON tracks strict request media negotiation for an NDJSON stream.
func BenchmarkNewFromContentTypeNDJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), http.MethodPost, "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.NDJSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromContentType(req)
	}
}

// BenchmarkNewFromContentTypeWithParameters tracks the parser path for a parameterized NDJSON request body.
func BenchmarkNewFromContentTypeWithParameters(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), http.MethodPost, "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.NDJSON+"; charset=utf-8")

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.StreamContent.NewFromContentType(req)
	}
}
