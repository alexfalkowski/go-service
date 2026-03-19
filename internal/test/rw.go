package test

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// BadReaderWriter is an os.ReaderWriter test double whose read and write operations fail.
type BadReaderWriter struct{}

// Exists reports true so callers continue into failing read and write paths.
func (rw *BadReaderWriter) Exists() bool {
	return true
}

// Read returns ErrFailed.
func (rw *BadReaderWriter) Read() ([]byte, error) {
	return nil, ErrFailed
}

// Write returns ErrFailed.
func (rw *BadReaderWriter) Write(_ []byte, _ os.FileMode) error {
	return ErrFailed
}

// Kind returns an empty kind string.
func (rw *BadReaderWriter) Kind() string {
	return strings.Empty
}
