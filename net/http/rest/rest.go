package rest

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/sync"
)

var (
	mux  *http.ServeMux
	cont *content.Content
	pool *sync.BufferPool
)

// Register for rpc.
func Register(mu *http.ServeMux, ct *content.Content, p *sync.BufferPool) {
	mux, cont, pool = mu, ct, p
}
