package unary_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

// benchmarkMedia prevents the compiler from eliminating media negotiation work.
var benchmarkMedia unary.Media

// BenchmarkNewFromMediaJSON tracks the exact media-type fast path used on hot request paths.
func BenchmarkNewFromMediaJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.UnaryContent.NewFromMedia(media.JSON)
	}
}

// BenchmarkNewFromRequestJSON tracks request header media negotiation overhead for a common JSON body.
func BenchmarkNewFromRequestJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.JSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.UnaryContent.NewFromRequest(req)
	}
}

// BenchmarkNewFromMediaWithParameters tracks the parser path needed for parameterized media types.
func BenchmarkNewFromMediaWithParameters(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.UnaryContent.NewFromMedia("application/json; profile=test")
	}
}

// BenchmarkNewFromRequestBodyJSON tracks the strict request-body resolution's switch-hit fast path
// (exact "application/json", no parameters), avoiding [net/http/media.Parse] on the hot request path.
func BenchmarkNewFromRequestBodyJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.JSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.UnaryContent.NewFromRequestBody(req)
	}
}

// BenchmarkNewFromRequestBodyWithParameters tracks the strict request-body resolution's parser branch,
// taken for any parameterized Content-Type such as "application/json; charset=utf-8".
func BenchmarkNewFromRequestBodyWithParameters(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.JSON+"; charset=utf-8")

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia, _ = test.UnaryContent.NewFromRequestBody(req)
	}
}
