package cache

import (
	"github.com/alexfalkowski/go-service/v2/compress/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

// maxSizeWriter enforces the cache encoded-size limit at the writer boundary.
// It prevents encoders from growing the intermediate buffer past max_size before
// compression gets its own chance to validate the payload.
type maxSizeWriter struct {
	writer  io.Writer
	max     int64
	written int64
}

func (w *maxSizeWriter) Write(data []byte) (int, error) {
	remaining := w.max - w.written
	if int64(len(data)) > remaining {
		if remaining <= 0 {
			return 0, errors.ErrTooLarge
		}

		n, err := w.writer.Write(data[:int(remaining)])
		w.written += int64(n)
		if err != nil {
			return n, err
		}

		return n, errors.ErrTooLarge
	}

	n, err := w.writer.Write(data)
	w.written += int64(n)

	return n, err
}
