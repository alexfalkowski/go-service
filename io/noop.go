package io

// NoopWriter is a writer that does nothing.
type NoopWriter struct{}

// Write implements the io.Writer interface.
func (nw *NoopWriter) Write(_ []byte) (int, error) {
	return 0, nil
}
