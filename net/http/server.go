package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/bytes"
	config "github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/time"
)

// NewServer constructs an HTTP server with common timeout defaults and supported protocol settings.
//
// Timeouts are derived from options first (if present) and fall back to [time.DefaultTimeout]:
//   - read_timeout
//   - write_timeout
//   - idle_timeout
//   - read_header_timeout
//
// Additional low-level server tuning may be provided through options using:
//   - max_header_bytes
//   - http2_max_concurrent_streams
//   - http2_max_receive_buffer_per_connection
//   - http2_max_receive_buffer_per_stream
//
// Protocols are configured via Protocols().
//
// Note: [opts.NonNegativeDuration] uses MustParseDuration under the hood; invalid or negative option
// values will panic at server construction time.
func NewServer(options config.Map, handler Handler) *Server {
	server := &http.Server{
		Handler:           handler,
		ReadTimeout:       options.NonNegativeDuration("read_timeout", time.DefaultTimeout).Duration(),
		WriteTimeout:      options.NonNegativeDuration("write_timeout", time.DefaultTimeout).Duration(),
		IdleTimeout:       options.NonNegativeDuration("idle_timeout", time.DefaultTimeout).Duration(),
		ReadHeaderTimeout: options.NonNegativeDuration("read_header_timeout", time.DefaultTimeout).Duration(),
		MaxHeaderBytes:    options.IntSize("max_header_bytes", bytes.Size(DefaultMaxHeaderBytes)),
		Protocols:         Protocols(),
	}
	http2 := &http.HTTP2Config{
		MaxConcurrentStreams:          int(options.Uint32("http2_max_concurrent_streams", 0)),
		MaxReceiveBufferPerConnection: int(options.Int32Size("http2_max_receive_buffer_per_connection", 0)),
		MaxReceiveBufferPerStream:     int(options.Int32Size("http2_max_receive_buffer_per_stream", 0)),
	}

	if http2.MaxConcurrentStreams != 0 || http2.MaxReceiveBufferPerConnection != 0 || http2.MaxReceiveBufferPerStream != 0 {
		server.HTTP2 = http2
	}

	return server
}
