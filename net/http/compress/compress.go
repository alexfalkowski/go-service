package compress

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/klauspost/compress/gzhttp"
)

// HeaderNoCompression is an alias for [gzhttp.HeaderNoCompression].
//
// Setting this response header before the response is written disables gzip compression for that response.
const HeaderNoCompression = gzhttp.HeaderNoCompression

// GzipHandler wraps h so its response is transparently gzip-compressed for clients that advertise support.
//
// This is a thin wrapper around [gzhttp.GzipHandler].
func GzipHandler(h http.Handler) http.HandlerFunc {
	return gzhttp.GzipHandler(h)
}

// Transport wraps rt so outbound requests advertise gzip support and gzipped responses are transparently decompressed.
//
// This is a thin wrapper around [gzhttp.Transport] with gzip compression enabled.
func Transport(rt http.RoundTripper) http.RoundTripper {
	return gzhttp.Transport(rt, gzhttp.TransportEnableGzip(true))
}
