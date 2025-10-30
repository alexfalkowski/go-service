package test

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// BadReaderWriter for test.
type BadReaderWriter struct{}

// Exists for test.
func (rw *BadReaderWriter) Exists() bool {
	return true
}

// Read for test.
func (rw *BadReaderWriter) Read() ([]byte, error) {
	return nil, ErrFailed
}

// Write for test.
func (rw *BadReaderWriter) Write(_ []byte, _ os.FileMode) error {
	return ErrFailed
}

// Kind for test.
func (rw *BadReaderWriter) Kind() string {
	return strings.Empty
}
