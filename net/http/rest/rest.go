package rest

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

var (
	mux  *http.ServeMux
	cont *content.Content
)

// Register for rpc.
func Register(mu *http.ServeMux, ct *content.Content) {
	mux, cont = mu, ct
}
