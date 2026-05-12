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

func BenchmarkNewFromMediaJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromMedia(media.JSON)
	}
}

func BenchmarkNewFromRequestJSON(b *testing.B) {
	req := httptest.NewRequestWithContext(b.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.JSON)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromRequest(req)
	}
}

func BenchmarkNewFromMediaWithCharset(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMedia = test.Content.NewFromMedia("application/json; charset=utf-8")
	}
}
