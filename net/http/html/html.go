package html

import (
	"net/http"
)

var mux *http.ServeMux

// Register for html.
func Register(mu *http.ServeMux) {
	mux = mu
}
