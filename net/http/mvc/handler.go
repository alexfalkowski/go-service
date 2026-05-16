package mvc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	snoop "github.com/felixge/httpsnoop"
)

// NotFoundHandler returns the MVC not-found handler for mixed HTTP transports.
//
// The handler only handles requests that explicitly accept HTML or come from HTMX, allowing API clients
// to receive the transport's content error fallback instead.
func NotFoundHandler() http.NotFoundHandler {
	return func(res http.ResponseWriter, req *http.Request) bool {
		if !IsDefined() || notFoundController == nil {
			return false
		}

		acceptsHTML := strings.Contains(req.Header.Get("Accept"), media.HTML)
		isHTMX := req.Header.Get("Hx-Request") == "true"
		if !acceptsHTML && !isHTMX {
			return false
		}

		writeNotFound(req, res)
		return true
	}
}

// NewHandler wraps mux with MVC error handling when MVC is defined.
//
// When MVC is not defined, it returns mux unchanged. When no NotFoundController has been registered, the
// returned handler preserves the mux's default behavior until one is registered.
func NewHandler(mux *http.ServeMux) http.Handler {
	if !IsDefined() {
		return mux
	}

	return &Handler{mux: mux}
}

// Handler wraps an HTTP mux and renders MVC error pages for unmatched routes.
type Handler struct {
	mux *http.ServeMux
}

// ServeHTTP serves req through the wrapped mux.
//
// If the mux has no matching route and the generated response is a 404, ServeHTTP renders the registered
// MVC error view instead of the mux's default plain-text response.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if !IsDefined() || notFoundController == nil {
		h.mux.ServeHTTP(res, req)
		return
	}

	handler, pattern := h.mux.Handler(req)
	if !strings.IsEmpty(pattern) {
		h.mux.ServeHTTP(res, req)
		return
	}

	buffer := pool.Get()
	defer pool.Put(buffer)

	response := &bufferedWriter{response: res, buffer: buffer, header: http.Header{}}
	metrics := snoop.CaptureMetricsFn(response, func(res http.ResponseWriter) {
		handler.ServeHTTP(res, req)
	})
	if metrics.Code != http.StatusNotFound {
		response.Flush()
		return
	}

	writeNotFound(req, res)
}
