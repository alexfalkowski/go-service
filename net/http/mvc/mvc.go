package mvc

import (
	"net/http"
)

var (
	mux   *http.ServeMux
	views *Views
)

// Register for mvc.
func Register(mu *http.ServeMux, vi *Views) {
	mux, views = mu, vi
}
