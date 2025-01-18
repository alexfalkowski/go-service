package test

// BadReaderCloser for test.
type BadReaderCloser struct{}

// Read returns ErrFailed.
func (r *BadReaderCloser) Read(_ []byte) (int, error) {
	return 0, ErrFailed
}

// Close returns ErrFailed.
func (r *BadReaderCloser) Close() error {
	return ErrFailed
}
