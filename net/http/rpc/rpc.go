package rpc

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
)

var (
	mux *http.ServeMux
	enc *encoding.Map
)

// Register for HTTP.
func Register(mu *http.ServeMux, en *encoding.Map) {
	mux, enc = mu, en
}
