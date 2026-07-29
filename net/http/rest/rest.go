package rest

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

var (
	router         *http.Router
	cont           *content.Content
	pool           *sync.BufferPool
	timeout        time.Duration
	maxReceiveSize bytes.Size
)

// Register stores the dependencies used by server and client helpers in package-level variables.
//
// Register is expected to be called during application startup (typically via dependency injection).
//
// Important: Register must be called before using any server-side route helpers or client helpers in this
// package. If it is not called, globals will be nil and helper calls will panic.
//
// After registration:
//   - server-side helpers (Get/Post/etc.) register handlers on the registered router, and
//   - client helpers (NewClient) build clients using the registered content codecs and buffer pool.
//
// t is the streaming inactivity timeout applied by this package's streaming route helpers (StreamRoute,
// StreamGet, StreamRouteRequest, StreamPost, StreamPut, StreamPatch). A successful Send extends the
// response write deadline by this duration. For bidirectional routes, a successful Recv also extends the
// request read deadline (see
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewStreamHandler] and
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestStreamHandler]), turning the
// whole-stream read and write timeouts into per-message inactivity budgets. Zero disables the extensions.
//
// m is the per-value request body size cap applied by this package's bidirectional streaming route
// helpers (StreamRouteRequest, StreamPost, StreamPut, StreamPatch). Unlike the buffered request-body
// path's cumulative limit, this bounds each value decoded by
// [github.com/alexfalkowski/go-service/v2/net/http/content.RequestStream.Recv] independently — a
// long-lived stream with many small values is never rejected for its cumulative size (see decision B18
// in the streaming design). Zero disables the per-value cap entirely.
func Register(r *http.Router, c *content.Content, p *sync.BufferPool, t time.Duration, m bytes.Size) {
	router = r
	cont = c
	pool = p
	timeout = t
	maxReceiveSize = m
}
