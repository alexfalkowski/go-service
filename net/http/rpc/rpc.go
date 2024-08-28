package rpc

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/sync"
)

var (
	mux  *http.ServeMux
	enc  *encoding.Map
	pool *sync.BufferPool
)

// Register for rpc.
func Register(mu *http.ServeMux, en *encoding.Map, p *sync.BufferPool) {
	mux, enc, pool = mu, en, p
}
