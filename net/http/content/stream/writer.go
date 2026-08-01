package stream

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// commitWriter buffers writes until commit is called, after which writes go straight to the live
// response writer.
//
// Buffering only the first value, rather than swapping the encoder to a new writer after commit, keeps
// one encoder instance bound to one writer for the life of the stream — required for codecs (yaml) that
// carry document-separator state tied to their writer.
type commitWriter struct {
	res       http.ResponseWriter
	buffer    *bytes.Buffer
	committed bool
}

// Write implements io.Writer. Before commit, writes accumulate in buffer; after commit, they go
// straight to res.
func (w *commitWriter) Write(p []byte) (int, error) {
	if w.committed {
		return w.res.Write(p)
	}

	return w.buffer.Write(p)
}

// commit flushes any buffered bytes to the live response writer. It is a no-op once already committed.
func (w *commitWriter) commit() error {
	if w.committed {
		return nil
	}

	w.committed = true
	_, err := w.buffer.WriteTo(w.res)

	return err
}
