package http

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-sync"
)

// bufferedWriter captures a mux-generated response until notFoundHandler decides whether to replay it.
//
// The not-found handler needs to inspect generated mux responses without committing the client response
// early: non-404 responses are flushed unchanged, while 404 responses can be replaced by a package-specific
// handler such as MVC not-found rendering or content status errors.
type bufferedWriter struct {
	header   Header
	response ResponseWriter
	buffer   *bytes.Buffer
	pool     *sync.BufferPool
	code     int
}

func (w *bufferedWriter) Header() Header {
	return w.header
}

func (w *bufferedWriter) Write(p []byte) (int, error) {
	if w.code == 0 {
		w.WriteHeader(StatusOK)
	}

	return w.buffer.Write(p)
}

func (w *bufferedWriter) WriteHeader(code int) {
	if w.code != 0 {
		return
	}

	w.code = code
}

func (w *bufferedWriter) Flush() {
	header := w.response.Header()
	for key, values := range w.header {
		if len(values) == 0 {
			header.Del(key)
			continue
		}

		header.Set(key, values[0])
		for _, value := range values[1:] {
			header.Add(key, value)
		}
	}

	if w.code != 0 {
		w.response.WriteHeader(w.code)
	}

	_, _ = w.buffer.WriteTo(w.response)
}

func (w *bufferedWriter) Close() {
	w.pool.Put(w.buffer)
}
