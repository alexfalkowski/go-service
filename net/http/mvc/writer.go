package mvc

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// bufferedWriter captures a mux-generated response until Handler decides whether to replay it.
//
// The MVC not-found handler needs to inspect generated mux responses without committing the client
// response early: non-404 responses are flushed unchanged, while 404 responses are replaced with the
// configured MVC not-found view.
type bufferedWriter struct {
	header   http.Header
	response http.ResponseWriter
	buffer   *bytes.Buffer
	code     int
}

func (w *bufferedWriter) Header() http.Header {
	return w.header
}

func (w *bufferedWriter) Write(p []byte) (int, error) {
	if w.code == 0 {
		w.WriteHeader(http.StatusOK)
	}

	return w.buffer.Write(p)
}

func (w *bufferedWriter) WriteHeader(code int) {
	w.code = code
}

func (w *bufferedWriter) Flush() {
	header := w.response.Header()
	for key, values := range w.header {
		header.Del(key)
		for _, value := range values {
			header.Add(key, value)
		}
	}

	if w.code != 0 {
		w.response.WriteHeader(w.code)
	}

	_, _ = w.buffer.WriteTo(w.response)
}
