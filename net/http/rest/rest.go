package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/sync"
)

var (
	mux  *http.ServeMux
	cont *content.Content
	pool *sync.BufferPool
)

// Register for rest.
func Register(mu *http.ServeMux, ct *content.Content, p *sync.BufferPool) {
	mux, cont, pool = mu, ct, p
}
