package rpc

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/sync"
)

var (
	mux  *http.ServeMux
	enc  *encoding.EncoderMap
	pool *sync.BufferPool
)

// Register for rpc.
func Register(mu *http.ServeMux, en *encoding.EncoderMap, p *sync.BufferPool) {
	mux, enc, pool = mu, en, p
}
