package mvc

import (
	"net/http"
)

var mux *http.ServeMux

// Register for mvc.
func Register(mu *http.ServeMux) {
	mux = mu
}
