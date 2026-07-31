package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-sync"
)

var (
	router *http.Router
	cont   *content.Content
	pool   *sync.BufferPool
	opts   content.StreamOptions
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
// so is passed through unchanged to this package's streaming route helpers (StreamRoute, StreamGet,
// StreamRouteRequest, StreamPost, StreamPut, StreamPatch) via
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewStreamHandler] and
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestStreamHandler]; see those
// constructors for the timeout and per-value receive-size semantics.
func Register(r *http.Router, c *content.Content, p *sync.BufferPool, so content.StreamOptions) {
	router = r
	cont = c
	pool = p
	opts = so
}
