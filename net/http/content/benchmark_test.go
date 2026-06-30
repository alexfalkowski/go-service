package content_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

// benchmarkMedia prevents the compiler from eliminating media negotiation work.
var benchmarkMedia content.Media

// BenchmarkNewFromMediaJSON tracks the exact media-type fast path used on hot request paths.
func BenchmarkNewFromMediaJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromMedia(media.JSON)
	}
}

// BenchmarkNewFromRequestJSON tracks request header media negotiation overhead for a common JSON body.
func BenchmarkNewFromRequestJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.JSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromRequest(req)
	}
}

// BenchmarkNewFromMediaWithParameters tracks the parser path needed for parameterized media types.
func BenchmarkNewFromMediaWithParameters(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromMedia("application/json; profile=test")
	}
}
